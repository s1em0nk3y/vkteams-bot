package message

import "io"

type Message struct {
	ChatID         string
	Text           string
	ReplyMsgID     string
	ForwardChatID  string
	ForwardMsgID   string
	KeyboardMarkup *KeyboardMarkup
	MessageFormat  *MessageFormat
	ParseMode      ParseMode
}

type FileMessage struct {
	Message
	FileID   string
	Filename string
	Contents io.Reader
}

type EditMessage struct {
	Message
	MessageToEditID string
}

type DeleteMessage struct {
	ChatID     string
	MessageIDs []string
}

type AnswerCallback struct {
	QueryID   string
	Text      string
	ShowAlert bool
	URL       string
}

type KeyboardMarkup [][]Button

type Button struct {
	Text     string      `json:"text"`
	URL      string      `json:"url,omitempty"`
	Callback string      `json:"callbackData,omitempty"`
	Style    ButtonStyle `json:"style,omitempty"`
}

// Button styles
type ButtonStyle string

const (
	ButtonBase      ButtonStyle = "base"
	ButtonPrimary   ButtonStyle = "primary"
	ButtonAttention ButtonStyle = "attention"
)

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
