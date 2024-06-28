package main

import (
	"fmt"
	"strconv"
	"strings"
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
	return fmt.Sprintf("https://t.me/c/%s/%d", convertChatId(chatId), msgId)
}

func text2Via(text string, chatId int64, msgId int64, fullname string) string {
	return fmt.Sprintf("%s\n[Via %s](%s)", text, escapeMarkdownV2(fullname), generateTelegramLink(chatId, msgId))
}

func convertChatId(chatId int64) string {
	text := strconv.FormatInt(chatId, 10)
	return strings.TrimPrefix(text, "-100")
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
