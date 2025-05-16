package event

type PartType string

const (
	PartTypeSticker PartType = "sticker"
	PartTypeMention PartType = "mention"
	PartTypeVoice   PartType = "voice"
	PartTypeFile    PartType = "file"
	PartTypeForward PartType = "forward"
	PartTypeReply   PartType = "reply"
)

type Part struct {
	Type    PartType    `json:"type"`
	Payload PartPayload `json:"payload"`
}

type PartPayload struct {
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	UserID    string      `json:"userId"`
	FileID    string      `json:"fileId"`
	Caption   string      `json:"caption"`
	Type      string      `json:"type"`
	Message   PartMessage `json:"message"`
}

type PartMessage struct {
	From      Contact `json:"from"`
	MsgID     string  `json:"msgId"`
	Text      string  `json:"text"`
	Timestamp int     `json:"timestamp"`
}
