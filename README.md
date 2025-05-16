# VK Teams API

[Specification](https://teams.vk.com/botapi/)

## Installation
```bash
go get github.com/s1em0nk3y/vkteams-bot@v0.0.1
```

## Usage
> You can look at [Example](./example/main.go)
### Create new bot
```Go
package main
import "github.com/s1em0nk3y/vkteams-bot"

func main() {
    bot := vkteams.New(
		config.Token, // Required parameter
		vkteams.WithApiURL(config.URL), // Custom API URL
		vkteams.WithHTTPClient(httpClient), // Custom HTTP Client (default is http.DefaultClient)
	)
}
```

### Send messages

#### Send Text
> Sends text message with three buttons
```Go
messageID, err := bot.SendText(context.Background(), &message.Message{
    ChatID:     "s1em0nk3y@ya.ru",
    Text:       "Some Text",
    ReplyMsgID: "message-id",
    KeyboardMarkup: &message.KeyboardMarkup{
        {{Text: "Button1", Callback: "callbackData1"}, {Text: "Button2", Callback: "callbackData2"}},
        {{Text: "Url Button", URL: "https://some.url"}},
    },
})
```
#### Send files
> Sends file with name Anyfile.file with content in bytes.Buffer
```Go
	messageID, fileID, err := bot.SendFile(context.Background(), &message.FileMessage{
		Message: message.Message{
			ChatID: "s1em0nk3y@ya.ru",
		},
		Filename: "Anyfile.file",
		Contents: bytes.NewBuffer([]byte("content of the file")),
	})
```
> Sends file via his file id (this works for sendVoice also)
```Go
	messageID, fileID, err = bot.SendFile(context.Background(), &message.FileMessage{
		Message: message.Message{
			ChatID: "s1em0nk3y@ya.ru",
		},
		FileID: "some file id",
	})
```
> Send voice message
```Go
    messageID, fileID, err = bot.SendVoice(context.Background(), &message.FileMessage{
		Message: message.Message{
			ChatID: "s1em0nk3y@ya.ru",
		},
		Filename: "voice.aac",
		Contents: bytes.NewBuffer([]byte("content of the file")),
	})
```
#### Edit message
```Go
    bot.EditMessage(context.Background(), &message.EditMessage{
		Message:         message.Message{ChatID: "s1em0nk3y@ya.ru"},
		MessageToEditID: "any message ID",
	})
```
#### Delete Messages
```Go
    bot.DeleteMessages(context.Background(), &message.DeleteMessage{
		ChatID:     "s1em0nk3y@ya.ru",
		MessageIDs: []string{"id1", "id2"},
	})
```
#### Answering Callback
```Go
    bot.AnswerCallback(context.Background(), &message.AnswerCallback{
		Text:      "Answer",
		QueryID:   "id of query",
		ShowAlert: false,
	})
```

### Listening events
```Go
	eventChannel := bot.UpdatesChannel(ctx)
	for event := range eventChannel {
		_, err := bot.SendText(ctx, &message.Message{
			ChatID: event.Chat.ID,
			Text:   "Text | <i>" + event.Text + "</i> | <b>Bold Text</b>",
			KeyboardMarkup: &message.KeyboardMarkup{
				{
					{
						Text:     "Some Text 1",
						Callback: "someCallback2",
					},
					{
						Text:     "Some <b>Text</b> 2",
						Callback: "someCallback2",
					},
				},
				{
					{
						Text: "Some URL",
						URL:  "https://test.example",
					},
				},
			},
			ReplyMsgID: event.MessageID,
			ParseMode:  message.ParseModeHTML,
		})
		log.Err(err).Msg("Send message")
	}
```