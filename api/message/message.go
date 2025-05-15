package message

//go:generate gotests -exported -template testify -w ./message.go
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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

// /messages/{sendFile,sendVoice} (Get/Post)
func (s *MessageService) SendFile(ctx context.Context, msg *vkteams.FileMessageRequest) (msgID string, fileID string, err error) {
	params := url.Values{
		"chatId":  {msg.ChatID},
		"caption": {msg.Text},
	}
	if msg.FileID != "" {
		params.Set("fileId", msg.FileID)
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

	response := struct {
		Id     string `json:"msgId"`
		FileID string `json:"fileId"`
		Ok     bool   `json:"Ok"`
	}{}

	// Send POST; Upload file
	if msg.Contents != nil {
		buffer := &bytes.Buffer{}
		part := multipart.NewWriter(buffer)
		fileWriter, err := part.CreateFormFile("file", msg.Filename)
		if err != nil {
			return "", "", fmt.Errorf("unable to create file writer: %w", err)
		}
		_, err = io.Copy(fileWriter, msg.Contents)
		if err != nil {
			return "", "", fmt.Errorf("unable to copy file contents: %w", err)
		}
		part.Close()
		req, err := s.client.PerformRequest(ctx, http.MethodPost, "/messages/sendFile", params, buffer)
		if err != nil {
			return "", "", err
		}
		req.Header.Set("Content-Type", part.FormDataContentType())
		resp, err := s.client.Do(req)
		if err != nil {
			return "", "", fmt.Errorf("unable to upload file: %w", err)
		}
		defer resp.Body.Close()
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return "", "", fmt.Errorf("unable to decode response: %w", err)
		}
		return response.Id, response.FileID, nil
	}

	req, err := s.client.PerformRequest(ctx, http.MethodPost, "/messages/sendFile", params, nil)
	if err != nil {
		return "", "", err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("unable to upload file: %w", err)
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", "", fmt.Errorf("unable to decode response: %w", err)
	}
	if !response.Ok {
		return "", "", errors.New("response status is not ok")
	}
	return response.Id, msg.FileID, nil
}

// // /messages/editText
// func (s *MessageService) EditMessage(msg *vkteams.EditMessageRequest) error {}

// // /messages/deleteMessage
// func (s *MessageService) DeleteMessages(*vkteams.DeleteMessageRequest) error {}

// // /messages/answerCallbackQuery
// func (s *MessageService) AnswerCallback(*vkteams.AnswerCallbackRequest) error {}
