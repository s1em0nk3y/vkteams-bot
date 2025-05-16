package message

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type Client interface {
	PerformRequest(ctx context.Context, method string, path string, params url.Values, body io.Reader) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)
}
