package vkteams

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var TestCfg = struct {
	Token     string `env:"VK_TOKEN,required"`
	URL       string `env:"VK_URL" envDefault:"https://myteam.mail.ru/bot/v1"`
	Proxy     bool   `env:"PROXY_ENABLE"`
	SSLVerify bool   `env:"SSL_VERIFY"`
}{}

var httpClient = &http.Client{
	Transport: http.DefaultTransport,
}

var testLogger zerolog.Logger
var testBot *Bot

func TestMain(m *testing.M) {
	godotenv.Load()
	testLogger = zerolog.New(zerolog.NewConsoleWriter())
	if err := env.Parse(&TestCfg); err != nil {
		testLogger.Fatal().Err(err).Send()
	}
	if !TestCfg.Proxy {
		httpClient.Transport.(*http.Transport).Proxy = nil
	}
	httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: TestCfg.SSLVerify,
	}
	testBot = New(TestCfg.Token,
		WithApiURL(TestCfg.URL),
		WithHTTPClient(httpClient),
	)
	m.Run()
}

func TestNew(t *testing.T) {
	type args struct {
		token string
		opts  []Option
	}
	tests := []struct {
		name string
		args args
		want *Bot
	}{
		{
			name: "All default",
			args: args{},
			want: &Bot{client: http.DefaultClient, apiUrl: defaultUrl},
		},
		{
			name: "Custom url",
			args: args{opts: []Option{WithApiURL("https://example.api/v1")}},
			want: &Bot{apiUrl: "https://example.api/v1", client: http.DefaultClient},
		},
		{
			name: "Custom http client",
			args: args{opts: []Option{WithHTTPClient(httpClient)}},
			want: &Bot{apiUrl: defaultUrl, client: httpClient},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, New(tt.args.token, tt.args.opts...))
		})
	}
}

// func TestBot_DoWithContext(t *testing.T) {
// 	type args struct {
// 		ctx    context.Context
// 		method string
// 		path   string
// 		params url.Values
// 		body   io.Reader
// 	}
// 	tests := []struct {
// 		name      string
// 		args      args
// 		want      *http.Response
// 		assertion assert.ErrorAssertionFunc
// 	}{
// 		{
// 			name: "Incorrect path",
// 			args: args{
// 				ctx:    context.Background(),
// 				method: "GET",
// 				path:   string([]byte{0x00, 0x01}),
// 				params: nil,
// 				body:   nil,
// 			},
// 			want: nil,
// 			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
// 				var target *url.Error
// 				if errors.As(err, &target) {
// 					return true
// 				}
// 				tt.Errorf("unexpected err: %s", err)
// 				return false
// 			},
// 		},
// 		{
// 			name: "Error building request",
// 			args: args{
// 				ctx:    context.Background(),
// 				method: "NOT EXIST",
// 			},
// 			want: nil,
// 			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
// 				target := "net/http: invalid method \"NOT EXIST\""
// 				var err2 = err
// 				for err2 != nil {
// 					if err2.Error() == target {
// 						return true
// 					}
// 					err2 = errors.Unwrap(err2)
// 				}

// 				tt.Errorf("unexpected err: %s", err.Error())
// 				return false
// 			},
// 		},
// 		{
// 			name: "Error getting response (context.Timeout)",
// 			args: args{
// 				ctx: func() context.Context {
// 					ctx, cancel := context.WithCancel(context.Background())
// 					cancel()
// 					return ctx
// 				}(),
// 				method: "GET",
// 				path:   "/self/get",
// 			},
// 			want: nil,
// 			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
// 				c := context.Canceled
// 				if errors.Is(err, c) {
// 					return true
// 				}
// 				tt.Errorf("got %s wanted %s", err.Error(), c.Error())
// 				return false
// 			},
// 		},
// 		{
// 			name: "Correct use",
// 			args: args{
// 				ctx:    context.Background(),
// 				method: "GET",
// 				path:   "/self/get",
// 			},
// 			want: &http.Response{StatusCode: 200},
// 			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
// 				if err != nil {
// 					tt.Errorf("error is not nil: %s", err.Error())
// 					return false
// 				}
// 				return true
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b := New(TestCfg.Token,
// 				WithApiURL(TestCfg.URL),
// 				WithHTTPClient(httpClient),
// 			)
// 			got, err := b.DoWithContext(tt.args.ctx, tt.args.method, tt.args.path, tt.args.params, nil, tt.args.body)
// 			tt.assertion(t, err)
// 			if tt.want != nil {
// 				require.Equal(t, tt.want.StatusCode, got.StatusCode)
// 			}
// 		})
// 	}
// }

func TestBot_PerformRequest(t *testing.T) {
	type args struct {
		ctx    context.Context
		method string
		path   string
		params url.Values
		body   io.Reader
	}
	tests := []struct {
		name      string
		args      args
		want      *http.Request
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Corrupt URL",
			args: args{
				ctx:  context.TODO(),
				path: string([]byte{0x00, 0x01}),
			},
			want: nil,
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				var targetErr *url.Error
				return assert.ErrorAs(tt, err, &targetErr)
			},
		},
		{
			name: "Corrupt http method (NON EXIST)",
			args: args{
				ctx:    context.TODO(),
				method: "NON EXIST",
			},
			want: nil,
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				target := "net/http: invalid method \"NON EXIST\""
				var err2 = err
				for err2 != nil {
					if err2.Error() == target {
						return true
					}
					err2 = errors.Unwrap(err2)
				}

				tt.Errorf("unexpected err: %s", err.Error())
				return false
			},
		},
		{
			name: "Valid usage",
			args: args{
				ctx:    context.Background(),
				method: "GET",
				path:   "/self/get",
			},
			want: func() *http.Request {
				req, _ := http.NewRequest("GET", TestCfg.URL+"/self/get", nil)
				req.URL.RawQuery = url.Values{
					"token": {TestCfg.Token},
				}.Encode()
				return req
			}(),
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testBot.PerformRequest(tt.args.ctx, tt.args.method, tt.args.path, tt.args.params, tt.args.body)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBot_Do(t *testing.T) {
	type args *http.Request
	tests := []struct {
		name      string
		args      args
		want      *http.Response
		wantOk    bool
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Correct usage",
			args: func() *http.Request {
				req, _ := testBot.PerformRequest(context.Background(), "GET", "/self/get", nil, nil)
				return req
			}(),
			want: &http.Response{
				StatusCode: http.StatusOK,
			},
			wantOk: true,
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(tt, err)
			},
		},
		{
			name: "No token",
			args: func() *http.Request {
				req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, testBot.apiUrl+"/self/get", nil)
				return req
			}(),
			want: &http.Response{
				StatusCode: http.StatusOK,
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(tt, err)
			},
		},
		{
			name: "Nil request",
			args: nil,
			want: nil,
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "no request provided")
			},
		},
		{
			name: "Canceled request",
			args: func() *http.Request {
				ctx, cancel := context.WithCancel(context.TODO())
				cancel()
				req, _ := testBot.PerformRequest(ctx, "GET", "/self/get", nil, nil)
				return req
			}(),
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, context.Canceled)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testBot.Do(tt.args)
			tt.assertion(t, err)
			if err == nil {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
				return
			}
			if tt.want == nil {
				return
			}

			assert.Equal(t, tt.want.StatusCode, got.StatusCode)
			response := struct {
				Ok bool `json:"Ok"`
			}{}
			assert.NoError(t, json.NewDecoder(got.Body).Decode(&response))
			assert.Equal(t, tt.wantOk, response.Ok)
		})
	}
}
