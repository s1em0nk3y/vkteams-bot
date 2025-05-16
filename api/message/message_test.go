package message_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/s1em0nk3y/vkteams-bot"
	"github.com/s1em0nk3y/vkteams-bot/api/message"
	"github.com/stretchr/testify/assert"
)

var TestCfg = struct {
	Token         string `env:"VK_TOKEN,required"`
	URL           string `env:"VK_URL" envDefault:"https://myteam.mail.ru/bot/v1"`
	Proxy         bool   `env:"PROXY_ENABLE"`
	SSLVerify     bool   `env:"SSL_VERIFY"`
	ChatID        string `env:"VK_CHAT_ID,required"`
	MessageID     string `env:"MESSAGE_ID,required"`
	FileID        string `env:"FILE_ID,required"`
	VoiceFilePath string `env:"VOICE_FILE,required"`
	VoiceFileID   string `env:"VOICE_FILE_ID,required"`
}{}

var httpClient = &http.Client{
	Transport: http.DefaultTransport,
}

var testLogger = zerolog.New(zerolog.NewConsoleWriter())

var testBot *vkteams.Bot

func TestMain(m *testing.M) {
	godotenv.Load("../../.env")
	if err := env.Parse(&TestCfg); err != nil {
		log.Fatal(err)
	}
	if !TestCfg.Proxy {
		httpClient.Transport.(*http.Transport).Proxy = nil
	}
	httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: TestCfg.SSLVerify,
	}
	testBot = vkteams.New(TestCfg.Token,
		vkteams.WithApiURL(TestCfg.URL),
		vkteams.WithHTTPClient(httpClient),
	)
	m.Run()
}

func TestMessageService_SendText(t *testing.T) {
	type args struct {
		ctx context.Context
		msg *message.Message
	}
	tests := []struct {
		name      string
		args      args
		wantMsgID string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Correct use; Forwarding message",
			args: args{
				ctx: testLogger.WithContext(context.Background()),
				msg: &message.Message{
					ChatID:        TestCfg.ChatID,
					Text:          "Some <b>Test</b> Text",
					ForwardChatID: TestCfg.ChatID,
					ForwardMsgID:  TestCfg.MessageID,
					ParseMode:     message.ParseModeHTML,
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				if err != nil {
					tt.Errorf("wanted nil got %s", err.Error())
					return false
				}
				return true
			},
		},
		{
			name: "Correct use; Replying message",
			args: args{
				ctx: testLogger.WithContext(context.Background()),
				msg: &message.Message{
					ChatID:     TestCfg.ChatID,
					Text:       "Some Test Text",
					ReplyMsgID: TestCfg.MessageID,
					ParseMode:  message.ParseModeMarkdown,
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				if err != nil {
					tt.Errorf("wanted nil got %s", err.Error())
					return false
				}
				return true
			},
		},
		{
			name: "Reply to nonexist message",
			args: args{
				ctx: testLogger.WithContext(context.Background()),
				msg: &message.Message{
					ChatID:     TestCfg.ChatID,
					Text:       "Some Test Text",
					ReplyMsgID: "NON EXIST REPLY ID",
					ParseMode:  message.ParseModeMarkdown,
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "response status is not ok")
			},
		},
		{
			name: "Context canceled",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(testLogger.WithContext(context.Background()))
					cancel()
					return ctx

				}(),
				msg: &message.Message{
					ChatID: TestCfg.ChatID,
					Text:   "Context canceled",
				},
			},
			wantMsgID: "",
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, context.Canceled)
			},
		},
		{
			name: "Correct usage; create buttons",
			args: args{
				ctx: testLogger.WithContext(context.Background()),
				msg: &message.Message{
					ChatID:    TestCfg.ChatID,
					Text:      "Some Test Text",
					ParseMode: message.ParseModeHTML,
					KeyboardMarkup: &message.KeyboardMarkup{
						{
							{
								Text:     "First",
								Style:    message.ButtonAttention,
								Callback: "somecallback1",
							},
							{
								Text:  "Second Button, With Url",
								URL:   "https://example.com",
								Style: message.ButtonBase,
							},
							{
								Text:     "Third, Primary style button",
								Callback: "somecallback3",
								Style:    message.ButtonPrimary,
							},
						},
					},
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				if err != nil {
					tt.Errorf("wanted nil got %s", err.Error())
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := message.New(testBot)
			gotMsgID, err := s.SendText(tt.args.ctx, tt.args.msg)
			tt.assertion(t, err)
			if err == nil {
				assert.NotEmpty(t, gotMsgID)
			}
		})
	}
}

func TestMessageService_SendFile(t *testing.T) {
	type args struct {
		ctx context.Context
		msg *message.FileMessage
	}
	tests := []struct {
		name       string
		args       args
		wantMsgID  string
		wantFileID string
		assertion  assert.ErrorAssertionFunc
		wantsOk    bool
	}{
		{
			name: "Send raw data",
			args: args{
				ctx: context.Background(),
				msg: &message.FileMessage{
					Message: message.Message{
						ChatID: TestCfg.ChatID,
						Text:   "Some Text",
					},
					Contents: bytes.NewBuffer([]byte("some text")),
					Filename: "Filename.txt",
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(tt, err)
			},
			wantsOk: true,
		},
		{
			name: "Send by file id",
			args: args{
				ctx: context.Background(),
				msg: &message.FileMessage{
					Message: message.Message{
						ChatID: TestCfg.ChatID,
						Text:   "Description",
					},
					FileID: TestCfg.FileID,
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(tt, err)
			},
			wantsOk: true,
		},
		{
			name: "Send by non exist file id",
			args: args{
				ctx: context.Background(),
				msg: &message.FileMessage{
					Message: message.Message{
						ChatID: TestCfg.ChatID,
						Text:   "Description",
					},
					FileID: "NON EXIST FILE ID",
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "response status is not ok")
			},
		},
		{
			name: "Send Voice file",
			args: args{
				ctx: context.Background(),
				msg: &message.FileMessage{
					Message: message.Message{
						ChatID: TestCfg.ChatID,
					},
					Contents: func() io.Reader {
						file, err := os.Open("../../" + TestCfg.VoiceFilePath)
						if err != nil {
							t.Fatalf("unable to open file: %s", TestCfg.VoiceFilePath)
						}
						defer file.Close()
						by := &bytes.Buffer{}
						io.Copy(by, file)
						return by
					}(),
					Filename: "filename.aac",
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(tt, err)
			},
			wantsOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := message.New(testBot)
			gotMsgID, gotFileID, err := s.SendFile(tt.args.ctx, tt.args.msg)
			if tt.assertion != nil {
				tt.assertion(t, err)
			}
			if tt.wantsOk {
				assert.NotEmpty(t, gotMsgID)
				assert.NotEmpty(t, gotFileID)
			}
		})
	}
}

func TestMessageService_SendVoice(t *testing.T) {
	type args struct {
		ctx context.Context
		msg *message.FileMessage
	}
	tests := []struct {
		name       string
		args       args
		wantMsgID  string
		wantFileID string
		assertion  assert.ErrorAssertionFunc
		wantsOk    bool
	}{
		{
			name: "Send Voice file",
			args: args{
				ctx: context.Background(),
				msg: &message.FileMessage{
					Message: message.Message{
						ChatID: TestCfg.ChatID,
					},
					Contents: func() io.Reader {
						file, err := os.Open("../../" + TestCfg.VoiceFilePath)
						if err != nil {
							t.Fatalf("unable to open file: %s", TestCfg.VoiceFilePath)
						}
						defer file.Close()
						by := &bytes.Buffer{}
						io.Copy(by, file)
						return by
					}(),
					Filename: "filename.aac",
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(tt, err)
			},
			wantsOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := message.New(testBot)
			gotMsgID, gotFileID, err := s.SendVoice(tt.args.ctx, tt.args.msg)
			if tt.assertion != nil {
				tt.assertion(t, err)
			}
			if tt.wantsOk {
				assert.NotEmpty(t, gotMsgID)
				assert.NotEmpty(t, gotFileID)
			}
		})
	}
}

func TestMessageService_EditMessage(t *testing.T) {
	type args struct {
		ctx context.Context
		msg *message.EditMessage
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Edit Non exist message",
			args: args{
				ctx: testLogger.WithContext(context.Background()),
				msg: &message.EditMessage{
					Message: message.Message{
						ChatID: TestCfg.ChatID,
						Text:   "EDITED MESSAGE",
					},
					MessageToEditID: "Non exist message ID",
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, message.ErrNotOk)
			},
		},
		{
			name: "Correct edit",
			args: args{
				ctx: context.Background(),
				msg: &message.EditMessage{
					Message: message.Message{
						ChatID: TestCfg.ChatID,
						Text:   "EDITED MESSAGE",
					},
					MessageToEditID: TestCfg.MessageID,
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(tt, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := message.New(testBot)
			tt.assertion(t, s.EditMessage(tt.args.ctx, tt.args.msg))
		})
	}
}

func TestMessageService_DeleteMessages(t *testing.T) {
	type args struct {
		ctx context.Context
		msg *message.DeleteMessage
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := message.New(testBot)
			tt.assertion(t, s.DeleteMessages(tt.args.ctx, tt.args.msg))
		})
	}
}

func TestMessageService_AnswerCallback(t *testing.T) {
	type args struct {
		ctx    context.Context
		answer *message.AnswerCallback
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := message.New(testBot)
			tt.assertion(t, s.AnswerCallback(tt.args.ctx, tt.args.answer))
		})
	}
}
