package event

type ChatType string
type EventType string

const (
	// Chat types
	ChatTypePrivate ChatType = "private"
	ChatTypeGroup   ChatType = "group"
	ChatTypeChannel ChatType = "channel"
	// Event types
	EventNewMessage      EventType = "newMessage"
	EventEditedMessage   EventType = "editedMessage"
	EventDeletedMessage  EventType = "deletedMessage"
	EventPinnedMessage   EventType = "pinnedMessage"
	EventUnpinnedMessage EventType = "unpinnedMessage"
	EventNewChatMembers  EventType = "newChatMembers"
	EventLeftChatMembers EventType = "leftChatMembers"
)

type Event struct {
	ID      int       `json:"eventId"`
	Type    EventType `json:"type"`
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

	// For events where members added/removed
	MembersLeft []Contact `json:"leftMembers"`
	MembersNew  []Contact `json:"newMembers"`
	AddedBy     Contact   `json:"addedBy"`
	RemovedBy   Contact   `json:"removedBy"`
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
