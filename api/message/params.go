package message

import (
	"encoding/json"
	"net/url"

	"github.com/s1em0nk3y/vkteams-bot"
)

func buildParams(msg *vkteams.Message) url.Values {
	params := url.Values{
		"chatId": {msg.ChatID},
	}
	if msg.ReplyMsgID != "" {
		params.Set("replyMsgId", msg.ReplyMsgID)
	}
	if msg.ForwardMsgID != "" {
		params.Set("forwardMsgId", msg.ForwardMsgID)
		params.Set("forwardChatId", msg.ForwardChatID)
	}
	if msg.KeyboardMarkup != nil {
		bytes, _ := json.Marshal(msg.KeyboardMarkup)
		params.Set("inlineKeyboardMarkup", string(bytes))
	}
	if msg.ParseMode != vkteams.ParseModeUnknown {
		params.Set("parseMode", msg.ParseMode.String())
	}
	return params
}
