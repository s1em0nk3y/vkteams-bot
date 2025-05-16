package vkteams

import "net/http"

type Option func(*Bot)

func WithHTTPClient(cli *http.Client) Option {
	return func(b *Bot) {
		b.client = cli
	}
}

func WithApiURL(baseUrl string) Option {
	return func(b *Bot) {
		b.apiUrl = baseUrl
	}
}

func WithPollSeconds(seconds uint) Option {
	return func(b *Bot) {
		b.pollSeconds = seconds
	}
}
