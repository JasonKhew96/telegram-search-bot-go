package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type SearchBot struct {
	config     *Config
	db         *Database
	bot        *gotgbot.Bot
	loc        *time.Location
	adminCache map[int64][]gotgbot.ChatMember
}

func StartBot(databaseFile, configFile string) {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Fatalln(err)
	}

	config, err := ParseConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	database, err := NewDatabase(databaseFile, false)
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

	if ok, err := bot.DeleteMyCommands(nil); err != nil {
		log.Println(err)
	} else if !ok {
		log.Println("DeleteMyCommands failed")
	} else {
		log.Println("DeleteMyCommands succeeded")
	}

	m := SearchBot{
		config:     config,
		db:         database,
		bot:        bot,
		loc:        loc,
		adminCache: make(map[int64][]gotgbot.ChatMember),
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println(err)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: 1,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	dispatcher.AddHandler(handlers.NewCommand("dlog", m.commandDeleteResponse).SetTriggers([]rune("/!")))
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
			AllowedUpdates: []string{"message", "edited_message", "inline_query", "chat_member"},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Bot started as %s\n", bot.User.Username)
	updater.Idle()
}

func (m *SearchBot) deleteMsg(chatId, msgId int64) func() {
	return func() {
		_, err := m.bot.DeleteMessage(chatId, msgId, nil)
		if err != nil {
			log.Println(err)
		}
	}
}

func (m *SearchBot) commandDeleteResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Type == "private" {
		return nil
	}
	if ctx.EffectiveSender.User == nil {
		return nil
	}

	chat, err := m.db.GetChat(ctx.EffectiveChat.Id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows || !chat.Enabled {
		return nil
	}

	re, err := regexp.Compile(`https://t\.me/c/(\d+)/(\d+)`)
	if err != nil {
		return err
	}
	matches := re.FindStringSubmatch(ctx.EffectiveMessage.GetText())
	if len(matches) != 3 {
		return nil
	}

	chatId, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return err
	}
	chatId = convert2BotChatId(chatId)

	if chatId != ctx.EffectiveChat.Id {
		return nil
	}

	chatPeer, err := m.db.GetChatPeerCount(ctx.EffectiveChat.Id, ctx.EffectiveSender.Id())
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows || chatPeer <= 0 {
		return nil
	}

	msgId, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return err
	}

	admins, err := m.GetChatAdministrators(chatId)
	if err != nil {
		return err
	}
	isAdmin := false
	for _, admin := range admins {
		if admin.GetUser().Id == ctx.EffectiveSender.Id() {
			isAdmin = true
			break
		}
	}

	msg, err := m.db.GetMessage(chatId, msgId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		msg, err := ctx.EffectiveMessage.Reply(b, "Message not found", nil)
		if err != nil {
			return err
		}
		time.AfterFunc(10*time.Second, m.deleteMsg(chatId, msg.MessageId))
		return nil
	}
	if msg.FromID != ctx.EffectiveSender.Id() && !isAdmin {
		msg, err := ctx.EffectiveMessage.Reply(b, "Unauthorized", nil)
		if err != nil {
			return err
		}
		time.AfterFunc(10*time.Second, m.deleteMsg(chatId, msg.MessageId))
		return nil
	}

	if err := m.db.DeleteMessage(chatId, msgId); err != nil {
		return err
	}

	deletedMessage, err := ctx.EffectiveMessage.Reply(b, "Message deleted", nil)
	if err != nil {
		return err
	}
	time.AfterFunc(10*time.Second, m.deleteMsg(chatId, deletedMessage.MessageId))

	return nil
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

func isAdmin(status string) bool {
	return status == "administrator" || status == "creator"
}

func isMember(status string) bool {
	return status == "member" || isAdmin(status)
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
	if !isAdmin(old.GetStatus()) && isAdmin(new.GetStatus()) {
		admins, err := m.GetChatAdministrators(chat.Id)
		if err != nil {
			log.Println(err)
		} else {
			newAdmins := []gotgbot.ChatMember{}
			for _, admin := range admins {
				if admin.GetUser().Id == new.GetUser().Id {
					continue
				}
				newAdmins = append(newAdmins, admin)
			}
			m.adminCache[chat.Id] = newAdmins
		}
	}
	if isAdmin(old.GetStatus()) && !isAdmin(new.GetStatus()) {
		admins, err := m.GetChatAdministrators(chat.Id)
		if err != nil {
			log.Println(err)
		}
		newAdmins := []gotgbot.ChatMember{}
		for _, admin := range admins {
			if admin.GetUser().Id == new.GetUser().Id {
				continue
			}
			newAdmins = append(newAdmins, admin)
		}
		m.adminCache[chat.Id] = newAdmins
	}
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
	return true
}

func (m *SearchBot) inlineQueryResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	count, err := m.db.GetChatPeersCount(ctx.InlineQuery.From.Id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
	}
	if err == sql.ErrNoRows || count <= 0 {
		_, err = ctx.InlineQuery.Answer(b, []gotgbot.InlineQueryResult{gotgbot.InlineQueryResultCachedSticker{
			Id:            "unauthorized_sticker",
			StickerFileId: "CAACAgUAAxkDAAEFBIhjffVfXIFyngE4vR2Zg_uDkDS41gACMAsAAoB48FdrYCP5TE3CEh4E",
		}}, &gotgbot.AnswerInlineQueryOpts{
			CacheTime:  300,
			IsPersonal: true,
		})
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
		_, err = ctx.InlineQuery.Answer(b, []gotgbot.InlineQueryResult{gotgbot.InlineQueryResultArticle{
			Id:    "info",
			Title: "No results found",
			InputMessageContent: gotgbot.InputTextMessageContent{
				MessageText: ".",
			},
			Description: "Please refine your search query and try again",
		}}, &gotgbot.AnswerInlineQueryOpts{
			CacheTime:  300,
			IsPersonal: true,
		})
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
		title, err := trimUnicodeAddEllipsis(mnp.Text, 64)
		if err != nil {
			log.Println(err)
			continue
		}
		expandableQuote, err := trimUnicodeAddEllipsis(mnp.Text, 2048)
		if err != nil {
			log.Println(err)
			continue
		}
		results = append(results, gotgbot.InlineQueryResultArticle{
			Id:          strconv.FormatInt(mnp.MSGID, 10),
			Title:       title,
			Description: fmt.Sprintf("%s %s@%s", mnp.Timestamp.In(m.loc).Format(time.DateTime), mnp.FullName, mnp.Title),
			InputMessageContent: gotgbot.InputTextMessageContent{
				MessageText: text2Via(text2ExpandableQuote(expandableQuote), mnp.Message.ChatID, mnp.MSGID, mnp.FullName),
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
		if ctx.EditedMessage.EditDate-ctx.EditedMessage.Date > 2*24*60*60 { // 48 hours
			return nil
		}
	}

	text := ctx.EffectiveMessage.GetText()
	return m.db.UpsertMessage(ctx.EffectiveChat.Id, ctx.EffectiveSender.Id(), ctx.EffectiveMessage.MessageId, text, ctx.EffectiveMessage.Date)
}

func (m *SearchBot) GetChatAdministrators(chatId int64) ([]gotgbot.ChatMember, error) {
	if admins, ok := m.adminCache[chatId]; ok {
		return admins, nil
	}
	admins, err := m.bot.GetChatAdministrators(chatId, nil)
	if err != nil {
		return nil, err
	}
	return admins, nil
}
