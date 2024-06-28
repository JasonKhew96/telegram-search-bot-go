package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type SearchBot struct {
	config *Config
	db     *Database
	bot    *gotgbot.Bot
	loc    *time.Location
}

func StartBot() {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Fatalln(err)
	}

	config, err := ParseConfig()
	if err != nil {
		log.Fatalln(err)
	}

	database, err := NewDatabase(false)
	if err != nil {
		log.Fatalln(err)
	}
	defer database.Close()

	requestOpts := &gotgbot.RequestOpts{
		Timeout: 60 * time.Second,
	}
	if config.CustomBotAPI != "" {
		requestOpts.APIURL = config.CustomBotAPI
	}
	bot, err := gotgbot.NewBot(config.BotToken, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			DefaultRequestOpts: requestOpts,
		},
		RequestOpts: requestOpts,
	})
	if err != nil {
		log.Fatalln(err)
	}

	m := SearchBot{
		config: config,
		db:     database,
		bot:    bot,
		loc:    loc,
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println(err)
			return ext.DispatcherActionNoop
		},
	})
	updater := ext.NewUpdater(dispatcher, nil)

	dispatcher.AddHandler(handlers.NewCommand("start", m.commandStartStopResponse).SetTriggers([]rune("/!")))
	dispatcher.AddHandler(handlers.NewCommand("stop", m.commandStartStopResponse).SetTriggers([]rune("/!")))
	dispatcher.AddHandler(handlers.NewChatMember(m.chatMemberRequest, m.chatMemberResponse))
	dispatcher.AddHandler(handlers.NewInlineQuery(m.inlineQueryRequest, m.inlineQueryResponse))
	dispatcher.AddHandler(handlers.NewMessage(m.newMessageRequest, m.newMessageResponse).SetAllowChannel(true).SetAllowEdited(true))

	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates:    config.DropPendingUpdate,
		EnableWebhookDeletion: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout:        60,
			AllowedUpdates: []string{"message", "edited_message", "channel_post", "edited_channel_post", "inline_query", "chat_member"},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Bot started as %s\n", bot.User.Username)
	updater.Idle()
}

func (m *SearchBot) commandStartStopResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Type == "private" {
		return nil
	}
	if ctx.EffectiveSender.User == nil {
		return nil
	}

	isStart := strings.HasPrefix(ctx.EffectiveMessage.GetText()[1:], "start")

	isEffectiveUserAdmin := false
	admins, err := ctx.EffectiveChat.GetAdministrators(b, nil)
	if err != nil {
		return err
	}
	for _, admin := range admins {
		if admin.GetUser().Id == ctx.EffectiveSender.Id() {
			isEffectiveUserAdmin = true
			break
		}
	}
	if !isEffectiveUserAdmin {
		return nil
	}

	chat, err := m.db.GetChat(ctx.EffectiveChat.Id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		if err := m.db.UpsertChat(ctx.EffectiveChat.Id, ctx.EffectiveChat.Title, isStart); err != nil {
			return err
		}
		text := "This chat is enabled"
		if !isStart {
			text = "This chat is disabled"
		}
		_, err = ctx.EffectiveMessage.Reply(b, text, nil)
		return err
	}

	if (!chat.Enabled && isStart) || (chat.Enabled && !isStart) {
		if err := m.db.UpdateChat(ctx.EffectiveChat.Id, ctx.EffectiveChat.Title, isStart); err != nil {
			return err
		}
	}

	var text string
	if chat.Enabled {
		if isStart {
			text = "This chat is already enabled"
		} else {
			text = "This chat is disabled"
		}
	} else {
		if !isStart {
			text = "This chat is already disabled"
		} else {
			text = "This chat is enabled"
		}
	}
	_, err = ctx.EffectiveMessage.Reply(b, text, nil)
	return err
}

func isMember(status string) bool {
	return status == "member" || status == "administrator" || status == "creator"
}

func isNotMember(status string) bool {
	return status == "left" || status == "kicked"
}

func (m *SearchBot) chatMemberRequest(u *gotgbot.ChatMemberUpdated) bool {
	chat, err := m.db.GetChat(u.Chat.Id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return false
	}
	if err == sql.ErrNoRows {
		return false
	}
	return chat.Enabled
}

func (m *SearchBot) chatMemberResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.ChatMember.Chat
	old := ctx.ChatMember.OldChatMember
	new := ctx.ChatMember.NewChatMember
	if isNotMember(old.GetStatus()) && isMember(new.GetStatus()) {
		// new member
		user := new.GetUser()
		if err := m.db.UpsertPeer(user.Id, strings.TrimSpace(fmt.Sprintf("%s %s", user.FirstName, user.LastName)), user.Username); err != nil {
			return err
		}
	}
	if isMember(old.GetStatus()) && isNotMember(new.GetStatus()) {
		// left or kicked
		if err := m.db.DeleteChatPeer(chat.Id, new.GetUser().Id); err != nil {
			return err
		}
	}
	return nil
}

func (m *SearchBot) inlineQueryRequest(iq *gotgbot.InlineQuery) bool {
	count, err := m.db.GetChatPeersCount(iq.From.Id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
	}
	if err == sql.ErrNoRows {
		return false
	}
	return count > 0
}

func (m *SearchBot) inlineQueryResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	count, err := m.db.GetChatPeersCount(ctx.InlineQuery.From.Id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
	}
	if err == sql.ErrNoRows || count <= 0 {
		_, err = b.AnswerInlineQuery(ctx.InlineQuery.Id, []gotgbot.InlineQueryResult{gotgbot.InlineQueryResultCachedSticker{
			Id:            "unauthorized_sticker",
			StickerFileId: "CAACAgUAAxkDAAEFBIhjffVfXIFyngE4vR2Zg_uDkDS41gACMAsAAoB48FdrYCP5TE3CEh4E",
		}}, nil)
		return err
	}

	chatPeers, err := m.db.GetChatPeersFromPeerId(ctx.InlineQuery.From.Id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return nil
	}
	var chatIds []int64
	for _, chatPeer := range chatPeers {
		chatIds = append(chatIds, chatPeer.ChatID)
	}

	// 2
	// text
	// text 2
	// @username text
	// @username text 2
	// @114514 text
	// @114514 text 2
	splits := strings.Split(ctx.InlineQuery.Query, " ")
	peerId := int64(0)
	username := ""
	queries := []string{}
	page := 1

	if strings.HasPrefix(splits[0], "@") {
		peerId, err = strconv.ParseInt(splits[0][1:], 10, 64)
		if err != nil {
			username = splits[0][1:]
		}
	}
	n, err := strconv.Atoi(splits[len(splits)-1])
	if err == nil && n > 1 {
		page = n
	}
	maxIndex := len(splits)
	if page > 1 {
		maxIndex -= 1
	}
	for i := 0; i < maxIndex; i++ {
		q := splits[i]
		if strings.HasPrefix(q, "@") {
			continue
		}
		if q == "" {
			continue
		}
		queries = append(queries, q)
	}

	messageAndPeers, err := m.db.SearchMessages(chatIds, username, peerId, queries, (page-1)*49)
	if err != nil {
		return err
	}

	if len(messageAndPeers) <= 0 {
		_, err = b.AnswerInlineQuery(ctx.InlineQuery.Id, []gotgbot.InlineQueryResult{gotgbot.InlineQueryResultArticle{
			Id:    "info",
			Title: "No results found",
			InputMessageContent: gotgbot.InputTextMessageContent{
				MessageText: ".",
			},
			Description: "Please refine your search query and try again",
		}}, nil)
		return err
	}

	results := []gotgbot.InlineQueryResult{
		gotgbot.InlineQueryResultArticle{
			Id:    "info",
			Title: fmt.Sprintf("Total %d Page: %d / %d", messageAndPeers[0].TotalCount, page, int64(math.Ceil(float64(messageAndPeers[0].TotalCount)/49))),
			InputMessageContent: gotgbot.InputTextMessageContent{
				MessageText: ".",
			},
		},
	}

	for _, mnp := range messageAndPeers {
		text := mnp.Text
		if len(text) > 100 {
			text = text[:100]
		}
		results = append(results, gotgbot.InlineQueryResultArticle{
			Id:          strconv.FormatInt(mnp.MSGID, 10),
			Title:       text,
			Description: fmt.Sprintf("%s %s@%s", mnp.Timestamp.In(m.loc).Format(time.DateTime), mnp.FullName, mnp.Title),
			InputMessageContent: gotgbot.InputTextMessageContent{
				MessageText: text2Via(text2ExpandableQuote(mnp.Text), mnp.Message.ChatID, mnp.MSGID, mnp.FullName),
				ParseMode:   "MarkdownV2",
				LinkPreviewOptions: &gotgbot.LinkPreviewOptions{
					IsDisabled: true,
				},
			},
		})
	}

	_, err = ctx.InlineQuery.Answer(b, results, &gotgbot.AnswerInlineQueryOpts{
		CacheTime:  15,
		IsPersonal: true,
	})
	return err
}

func (m *SearchBot) newMessageRequest(msg *gotgbot.Message) bool {
	if msg.Chat.Type == "private" {
		return false
	}
	if msg.From != nil && msg.From.IsBot {
		return false
	}
	if msg.ViaBot != nil && msg.ViaBot.Id == m.bot.Id {
		return false
	}
	if msg.GetText() == "" {
		return false
	}
	chat, err := m.db.GetChat(msg.Chat.Id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return false
	}
	if err == sql.ErrNoRows {
		return false
	}
	return chat.Enabled
}

func (m *SearchBot) newMessageResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	if err := m.db.UpsertPeer(ctx.EffectiveSender.Id(), strings.TrimSpace(fmt.Sprintf("%s %s", ctx.EffectiveSender.FirstName(), ctx.EffectiveSender.LastName())), ctx.EffectiveSender.Username()); err != nil {
		return err
	}
	if err := m.db.InsertChatPeer(ctx.EffectiveChat.Id, ctx.EffectiveSender.Id()); err != nil {
		return err
	}

	if ctx.EditedMessage != nil {
		if ctx.EditedMessage.EditDate-ctx.EditedMessage.Date > 120 {
			return nil
		}
	}

	text := ctx.EffectiveMessage.GetText()
	return m.db.UpsertMessage(ctx.EffectiveChat.Id, ctx.EffectiveSender.Id(), ctx.EffectiveMessage.MessageId, text, ctx.EffectiveMessage.Date)
}
