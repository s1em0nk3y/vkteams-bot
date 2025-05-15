package event

type ChatType string

const (
	ChatTypePrivate = "private"
	ChatTypeGroup   = "group"
	ChatTypeChannel = "channel"
)

type Event struct {
	ID      int    `json:"eventId"`
	Type    string `json:"type"`
	Payload `json:"payload"`
}

type Payload struct {
	BasePayload
	// Parts of message (sticker, file etc.)
	Parts []Part `json:"parts"`

	// For callback
	QueryID         string      `json:"queryId"`
	CallbackMessage BasePayload `json:"message"`
	CallbackData    string      `json:"callbackData"`
}

type BasePayload struct {
	MessageID string  `json:"msgId"`
	Chat      Chat    `json:"chat"`
	From      Contact `json:"from"`
	Timestamp int     `json:"timestamp"`
	Text      string  `json:"text"`
	// TODO: Format `json:"format"`
	EditedTimestamp int `json:"editedTimestamp"`
}

type Contact struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Chat struct {
	ID    string   `json:"chatId"`
	Type  ChatType `json:"type"`
	Title string   `json:"title"`
}
