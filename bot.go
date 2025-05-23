package vkteams

//go:generate gotests -exported -template testify -w ./bot.go
import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog"
	"github.com/s1em0nk3y/vkteams-bot/api/event"
	"github.com/s1em0nk3y/vkteams-bot/api/message"
)

const defaultUrl = "https://myteam.mail.ru/bot/v1"

type Bot struct {
	client      *http.Client
	apiUrl      string
	token       string
	pollSeconds uint
	*message.MessageService
	*event.EventService
}

func New(token string, opts ...Option) *Bot {
	b := &Bot{
		client:      http.DefaultClient,
		apiUrl:      defaultUrl,
		token:       token,
		pollSeconds: 60,
	}
	for _, opt := range opts {
		opt(b)
	}
	b.EventService = event.New(b, b.pollSeconds)
	b.MessageService = message.New(b)
	return b
}

func (b *Bot) PerformRequest(ctx context.Context, method string, path string, params url.Values, body io.Reader) (*http.Request, error) {
	log := *zerolog.Ctx(ctx)
	log = log.With().Str("path", b.apiUrl+path).Logger()
	urlPath, err := url.Parse(b.apiUrl + path)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url: %w", err)
	}
	if params == nil {
		params = url.Values{}
	}
	params.Set("token", b.token)
	urlPath.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, method, urlPath.String(), body)
	log.Err(err).Msg("create request")
	return req, err
}

func (b *Bot) Do(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("no request provided")
	}
	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error occured when sending request: %w", err)
	}
	return resp, err
}
