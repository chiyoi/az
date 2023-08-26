package graph

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/chiyoi/apricot/kitsune"
)

const (
	Endpoint = "https://graph.microsoft.com/"
)

func Read(ctx context.Context, token string, v any, query ...string) (err error) {
	re, err := Request(ctx, http.MethodGet, nil, token, query...)
	if err != nil {
		return
	}
	defer re.Body.Close()

	return kitsune.ParseResponse(re, v)
}

func Request(ctx context.Context, method string, body io.Reader, token string, query ...string) (re *http.Response, err error) {
	u, err := url.JoinPath(Endpoint, query...)
	if err != nil {
		return
	}

	r, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return
	}
	kitsune.SetAuthorization(r.Header, token)

	return http.DefaultClient.Do(r)
}
