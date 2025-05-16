package main

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/s1em0nk3y/vkteams-bot"
	"github.com/s1em0nk3y/vkteams-bot/api/message"
)

var config = struct {
	Token string `env:"VK_TOKEN,required"`
	URL   string `env:"VK_URL"`
	HTTP  struct {
		Proxy     bool `env:"PROXY"`
		SSLVerify bool `env:"SSL_VERIFY"`
	}
}{}

func main() {
	log := zerolog.New(zerolog.NewConsoleWriter())
	godotenv.Load()
	if err := env.Parse(&config); err != nil {
		log.Fatal().Err(err).Send()
	}
	// Configure HTTP client
	httpClient := &http.Client{
		Transport: http.DefaultTransport,
	}
	// Disable proxy (from env, if needed)
	if !config.HTTP.Proxy {
		httpClient.Transport.(*http.Transport).Proxy = nil
	}
	// Disable SSL verification if needed
	httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: config.HTTP.SSLVerify,
	}

	// Create bot instance
	bot := vkteams.New(
		config.Token,
		vkteams.WithApiURL(config.URL),
		vkteams.WithHTTPClient(httpClient),
	)
	ctx := log.WithContext(context.Background())

	// Listen Events
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
}
