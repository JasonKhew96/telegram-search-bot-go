package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/clipperhouse/uax29/graphemes"
)

var allMdV2 = []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
var mdV2Repl = strings.NewReplacer(func() (out []string) {
	for _, x := range allMdV2 {
		out = append(out, x, "\\"+x)
	}
	return out
}()...)

func escapeMarkdownV2(s string) string {
	return mdV2Repl.Replace(s)
}

func text2ExpandableQuote(text string) string {
	result := ""
	splits := strings.Split(text, "\n")
	for i, s := range splits {
		if i == 0 {
			result += "**>" + escapeMarkdownV2(s) + "\n"
		} else if i == len(splits)-1 {
			result += ">" + escapeMarkdownV2(s) + "||"
		} else {
			result += ">" + escapeMarkdownV2(s) + "\n"
		}
	}
	return result
}

func generateTelegramLink(chatId, msgId int64) string {
	return fmt.Sprintf("https://t.me/c/%s/%d", convert2NativeChatId(chatId), msgId)
}

func text2Via(text string, chatId int64, msgId int64, fullname string) string {
	return fmt.Sprintf("%s\n[Via %s](%s)", text, escapeMarkdownV2(fullname), generateTelegramLink(chatId, msgId))
}

func convert2NativeChatId(chatId int64) string {
	text := strconv.FormatInt(chatId, 10)
	return strings.TrimPrefix(text, "-100")
}

func convert2BotChatId(chatId int64) int64 {
	text := strconv.FormatInt(chatId, 10)
	newChatId, err := strconv.ParseInt("-100"+text, 10, 64)
	if err != nil {
		return chatId
	}
	return newChatId
}

func removeDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func trimUnicodeAddEllipsis(text string, l int) (string, error) {
	var newByteText []byte
	byteText := []byte(text)
	segments := graphemes.NewSegmenter(byteText)
	for segments.Next() && len(newByteText) <= l {
		newByteText = append(newByteText, segments.Bytes()...)
	}
	if err := segments.Err(); err != nil {
		return "", err
	}
	newText := string(newByteText)
	if len(text) > len(newText) {
		return newText + "…", nil
	}
	return newText, nil
}
