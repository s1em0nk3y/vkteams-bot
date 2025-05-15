package vkteams

import "context"

type MessageService interface {
	SendText(ctx context.Context, msg *Message) (msgID string, err error)
	SendFile(ctx context.Context, msg *FileMessage) (msgID string, fileID string, err error)
	SendVoice(ctx context.Context, msg *FileMessage) (msgID string, fileID string, err error)
	EditMessage(ctx context.Context, msg *EditMessage) error
	DeleteMessages(ctx context.Context, msg *DeleteMessage) error
	AnswerCallback(ctx context.Context, answer *AnswerCallback) error
}
