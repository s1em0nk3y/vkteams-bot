package vkteams

import "io"

type MessageRequest struct {
	ChatID         string
	Text           string
	ReplyMsgID     string
	ForwardChatID  string
	ForwardMsgID   string
	KeyboardMarkup *KeyboardMarkup
	MessageFormat  *MessageFormat
	ParseMode      ParseMode
}

type FileMessageRequest struct {
	MessageRequest
	FileID   string
	Filename string
	Contents io.Reader
}

type EditMessageRequest struct {
	ChatID         string
	Text           string
	MessageID      string
	KeyboardMarkup *KeyboardMarkup
	MessageFormat  *MessageFormat
	ParseMode      ParseMode
}

type DeleteMessageRequest struct {
	ChatID     string
	MessageIDs []string
}

type AnswerCallbackRequest struct {
	QueryID   string
	Text      string
	ShowAlert bool
	URL       string
}

type KeyboardMarkup struct{}

type MessageFormat struct{}

type ParseMode int

const (
	ParseModeUnknown ParseMode = iota
	ParseModeMarkdown
	ParseModeHTML
)

func (p ParseMode) String() string {
	switch p {
	case ParseModeHTML:
		return "HTML"
	case ParseModeMarkdown:
		return "MarkdownV2"
	default:
		return ""
	}
}
