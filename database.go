package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/JasonKhew96/telegram-search-bot-go/models"
	"github.com/liuzl/gocc"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	_ "modernc.org/sqlite"
)

/*
CREATE TABLE "chat" (
    "id" INTEGER NOT NULL,
    "title" TEXT NOT NULL,
    "enabled" BOOLEAN NOT NULL,
    PRIMARY KEY("id")
);

CREATE TABLE "chat_peer" (
    "id" INTEGER NOT NULL,
    "chat_id" INTEGER NOT NULL,
    "peer_id" INTEGER NOT NULL,
    PRIMARY KEY("id")
);

CREATE TABLE "message" (
    "id" TEXT NOT NULL,
    "chat_id" INTEGER NOT NULL,
    "from_id" INTEGER NOT NULL,
    "msg_id" INTEGER NOT NULL,
    "text" TEXT NOT NULL,
    "timestamp" DATETIME NOT NULL,
    "deleted" BOOLEAN NOT NULL DEFAULT 0,
    PRIMARY KEY("id")
);

CREATE INDEX "idx_message" ON "message" (
    "chat_id",
    "from_id",
    "msg_id",
    "text",
    "timestamp",
    "deleted"
);

CREATE TABLE "peer" (
    "id" INTEGER NOT NULL,
    "full_name" TEXT NOT NULL,
    "username" TEXT NOT NULL,
    PRIMARY KEY("id")
);
*/

type Database struct {
	db  *sql.DB
	ctx context.Context
	s2t *gocc.OpenCC
	t2s *gocc.OpenCC
}

type MessageAndPeer struct {
	TotalCount     int `boil:"total_count"`
	models.Message `boil:",bind"`
	models.Peer    `boil:",bind"`
	models.Chat    `boil:",bind"`
}

func NewDatabase(databaseFile string, importMode bool) (*Database, error) {
	// boil.DebugMode = true
	db, err := sql.Open("sqlite", fmt.Sprintf("%s?cache=shared", databaseFile))
	if err != nil {
		return nil, err
	}

	// migrations start
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "1_deleted_message",
				Up: []string{
					`ALTER TABLE "message" ADD COLUMN "deleted" BOOLEAN NOT NULL DEFAULT 0;`,
					`DROP INDEX IF EXISTS "idx_message_chat_id_msg_id";`,
					`CREATE INDEX "idx_message" ON "message" ("chat_id", "from_id", "msg_id", "text", "timestamp", "deleted");`,
				},
				Down: []string{
					`ALTER TABLE "message" DROP COLUMN "deleted";`,
					`DROP INDEX IF EXISTS "idx_message";`,
					`CREATE INDEX "idx_message_chat_id_msg_id" ON "message" ("chat_id", "msg_id", "text");`,
				},
			},
		},
	}
	migrationCount, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return nil, err
	}
	log.Printf("Applied %d migrations", migrationCount)
	// migrations end

	s2t, err := gocc.New("s2t")
	if err != nil {
		return nil, err
	}
	t2s, err := gocc.New("t2s")
	if err != nil {
		return nil, err
	}

	if importMode {
		// https://avi.im/blag/2021/fast-sqlite-inserts/
		db.Exec("PRAGMA journal_mode = OFF;")
		db.Exec("PRAGMA synchronous = 0;")
		db.Exec("PRAGMA cache_size = 1000000;")
		db.Exec("PRAGMA locking_mode = EXCLUSIVE;")
		db.Exec("PRAGMA temp_store = MEMORY;")
	}

	return &Database{
		db:  db,
		ctx: context.Background(),
		s2t: s2t,
		t2s: t2s,
	}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetChat(chatId int64) (*models.Chat, error) {
	return models.Chats(models.ChatWhere.ID.EQ(chatId)).One(d.ctx, d.db)
}

func (d *Database) UpdateChat(chatId int64, title string, enabled bool) error {
	chat, err := d.GetChat(chatId)
	if err != nil {
		return err
	}
	chat.Title = title
	chat.Enabled = enabled
	_, err = chat.Update(d.ctx, d.db, boil.Infer())
	return err
}

func (d *Database) UpsertChat(chatId int64, title string, enabled bool) error {
	chat := models.Chat{
		ID:      chatId,
		Title:   title,
		Enabled: enabled,
	}
	return chat.Upsert(d.ctx, d.db, true, []string{"id"}, boil.Infer(), boil.Infer())
}

func (d *Database) GetPeer(peerId int64) (*models.Peer, error) {
	return models.Peers(models.PeerWhere.ID.EQ(peerId)).One(d.ctx, d.db)
}

func (d *Database) UpsertPeer(peerId int64, fullName, username string) error {
	peer := models.Peer{
		ID:       peerId,
		FullName: fullName,
		Username: username,
	}
	return peer.Upsert(d.ctx, d.db, true, []string{"id"}, boil.Infer(), boil.Infer())
}

func (d *Database) GetMessageCount() (int64, error) {
	return models.Messages().Count(d.ctx, d.db)
}

func (d *Database) SearchMessages(chatId []int64, username string, peerId int64, texts []string, offset int) ([]*MessageAndPeer, error) {
	queryMods := []qm.QueryMod{qm.Select("message.msg_id", "message.chat_id", "message.text", "message.timestamp", "peer.full_name", "chat.title", "COUNT() OVER() as total_count"), qm.From("message"), qm.InnerJoin("peer on peer.id = message.from_id"), qm.InnerJoin("chat on chat.id = message.chat_id"), models.MessageWhere.Deleted.EQ(false), qm.Offset(offset), qm.Limit(49), qm.OrderBy("message.timestamp DESC")}
	for _, c := range chatId {
		queryMods = append(queryMods, models.MessageWhere.ChatID.EQ(c))
	}
	if username != "" {
		queryMods = append(queryMods, models.PeerWhere.Username.EQ(username))
	}
	if peerId != 0 {
		queryMods = append(queryMods, models.MessageWhere.FromID.EQ(peerId))
	}
	for _, q := range texts {
		t, err := d.s2t.Convert(q)
		if err != nil {
			log.Println(err)
		}
		s, err := d.t2s.Convert(q)
		if err != nil {
			log.Println(err)
		}
		dups := removeDuplicate([]string{q, t, s})
		var rawQuery string
		var rawArgs []interface{}
		if len(dups) == 1 {
			rawQuery += "message.text LIKE ?"
			rawArgs = append(rawArgs, "%"+dups[0]+"%")
		} else if len(dups) > 1 {
			for i, d := range dups {
				rawQuery += "message.text LIKE ?"
				rawArgs = append(rawArgs, "%"+d+"%")
				if i != len(dups)-1 {
					rawQuery += " OR "
				}
			}
		}
		queryMods = append(queryMods, qm.And(rawQuery, rawArgs...))
	}
	var messageAndPeer []*MessageAndPeer
	if err := models.NewQuery(queryMods...).Bind(d.ctx, d.db, &messageAndPeer); err != nil {
		return nil, err
	}
	return messageAndPeer, nil
}

func (d *Database) UpsertMessage(chatId int64, fromId int64, msgId int64, text string, timestamp int64) error {
	message := models.Message{
		ID:        strconv.FormatInt(chatId, 10) + "_" + strconv.FormatInt(msgId, 10),
		ChatID:    chatId,
		FromID:    fromId,
		MSGID:     msgId,
		Text:      text,
		Timestamp: time.Unix(timestamp, 0),
	}
	return message.Upsert(d.ctx, d.db, true, []string{"id"}, boil.Blacklist("deleted"), boil.Infer())
}

func (d *Database) GetChatPeersCount(peerId int64) (int64, error) {
	return models.ChatPeers(models.ChatPeerWhere.PeerID.EQ(peerId)).Count(d.ctx, d.db)
}

func (d *Database) GetChatPeersFromPeerId(peerId int64) ([]*models.ChatPeer, error) {
	return models.ChatPeers(models.ChatPeerWhere.PeerID.EQ(peerId)).All(d.ctx, d.db)
}

func (d *Database) GetChatPeers(chatId int64, peerId int64) (*models.ChatPeer, error) {
	return models.ChatPeers(models.ChatPeerWhere.ChatID.EQ(chatId), models.ChatPeerWhere.PeerID.EQ(peerId)).One(d.ctx, d.db)
}

func (d *Database) InsertChatPeer(chatId int64, peerId int64) error {
	_, err := d.GetChatPeers(chatId, peerId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		return nil
	}
	chatPeer := models.ChatPeer{
		ChatID: chatId,
		PeerID: peerId,
	}
	return chatPeer.Insert(d.ctx, d.db, boil.Infer())
}

func (d *Database) DeleteChatPeer(chatId int64, peerId int64) error {
	_, err := models.ChatPeers(models.ChatPeerWhere.ChatID.EQ(chatId), models.ChatPeerWhere.PeerID.EQ(peerId)).DeleteAll(d.ctx, d.db)
	return err
}
