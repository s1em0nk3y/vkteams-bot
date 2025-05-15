package message

//go:generate gotests -exported -template testify -w ./message.go
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rs/zerolog"
	"github.com/s1em0nk3y/vkteams-bot"
)

type MessageService struct {
	client *vkteams.Bot
}

// /messages/sendText (Get)
func (s *MessageService) SendText(ctx context.Context, msg *vkteams.MessageRequest) (msgID string, err error) {
	params := url.Values{
		"chatId": {msg.ChatID},
		"text":   {msg.Text},
	}
	log := zerolog.Ctx(ctx)
	if msg.ReplyMsgID != "" {
		params.Set("replyMsgId", msg.ReplyMsgID)
	}
	if msg.ForwardMsgID != "" {
		params.Set("forwardMsgId", msg.ForwardMsgID)
		params.Set("forwardChatId", msg.ForwardChatID)
	}
	if msg.KeyboardMarkup != nil {
		log.Error().Msg("Keyboard Markup not implemented")
	}
	if msg.ParseMode != vkteams.ParseModeUnknown {
		params.Set("parseMode", msg.ParseMode.String())
	}

	req, err := s.client.PerformRequest(ctx, http.MethodGet, "/messages/sendText", params, nil)
	if err != nil {
		return "", err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to send text: %w", err)
	}
	defer resp.Body.Close()

	response := struct {
		Id string `json:"msgId"`
		Ok bool   `json:"Ok"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("unable to decode response: %w", err)
	}
	if !response.Ok {
		return "", errors.New("response status is not ok")
	}
	return response.Id, nil
}

// /messages/sendFile (Get/Post)
func (s *MessageService) SendFile(ctx context.Context, msg *vkteams.FileMessageRequest) (msgID string, fileID string, err error) {
	return s.sendFile(ctx, msg, "/messages/sendFile")
}

func (s *MessageService) SendVoice(ctx context.Context, msg *vkteams.FileMessageRequest) (msgID string, fileID string, err error) {
	return s.sendFile(ctx, msg, "/messages/sendVoice")
}

// // /messages/editText
// func (s *MessageService) EditMessage(msg *vkteams.EditMessageRequest) error {}

// // /messages/deleteMessage
// func (s *MessageService) DeleteMessages(*vkteams.DeleteMessageRequest) error {}

// // /messages/answerCallbackQuery
// func (s *MessageService) AnswerCallback(*vkteams.AnswerCallbackRequest) error {}
