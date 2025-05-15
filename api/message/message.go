package message

//go:generate gotests -exported -template testify -w ./message.go
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/s1em0nk3y/vkteams-bot"
)

type MessageService struct {
	client *vkteams.Bot
}

// /messages/sendText (Get)
func (s *MessageService) SendText(ctx context.Context, msg *vkteams.Message) (msgID string, err error) {
	params := buildParams(msg)
	params.Set("text", msg.Text)
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
		return "", ErrNotOk
	}
	return response.Id, nil
}

// /messages/sendFile (Get/Post)
func (s *MessageService) SendFile(ctx context.Context, msg *vkteams.FileMessage) (msgID string, fileID string, err error) {
	return s.sendFile(ctx, msg, "/messages/sendFile")
}

// /messages/sendVoice (Get/Post)
func (s *MessageService) SendVoice(ctx context.Context, msg *vkteams.FileMessage) (msgID string, fileID string, err error) {
	return s.sendFile(ctx, msg, "/messages/sendVoice")
}

// /messages/editText
func (s *MessageService) EditMessage(ctx context.Context, msg *vkteams.EditMessage) error {
	params := buildParams(&msg.Message)
	params.Set("msgId", msg.MessageToEditID)
	params.Set("text", msg.Text)
	req, err := s.client.PerformRequest(ctx, http.MethodGet, "/messages/editText", params, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	response := struct {
		Ok bool `json:"Ok"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("unable to decode response: %w", err)
	}
	if !response.Ok {
		return ErrNotOk
	}
	return nil
}

// /messages/deleteMessage
func (s *MessageService) DeleteMessages(ctx context.Context, msg *vkteams.DeleteMessage) error {
	params := url.Values{
		"chatId": {msg.ChatID},
	}
	for _, msg := range msg.MessageIDs {
		params.Add("msgId", msg)
	}
	req, err := s.client.PerformRequest(ctx, http.MethodGet, "/messages/deleteMessages", params, nil)
	if err != nil {
		return err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	response := struct {
		Ok bool `json:"Ok"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("unable to decode response: %w", err)
	}
	if !response.Ok {
		return ErrNotOk
	}
	return nil
}

// /messages/answerCallbackQuery
func (s *MessageService) AnswerCallback(ctx context.Context, answer *vkteams.AnswerCallback) error {
	params := url.Values{
		"queryId": {answer.QueryID},
	}
	if answer.Text != "" {
		params.Set("text", answer.Text)
	}
	if answer.ShowAlert {
		params.Set("showAlert", "true")
	}
	if answer.URL != "" {
		params.Set("url", answer.URL)
	}

	req, err := s.client.PerformRequest(ctx, http.MethodGet, "/messages/answerCallbackQuery", params, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	response := struct {
		Ok bool `json:"Ok"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("unable to decode response: %w", err)
	}
	if !response.Ok {
		return ErrNotOk
	}
	return nil
}
