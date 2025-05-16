package message

import (
	"encoding/json"
	"net/url"
)

func buildParams(msg *Message) url.Values {
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
	if msg.ParseMode != ParseModeUnknown {
		params.Set("parseMode", msg.ParseMode.String())
	}
	return params
}
