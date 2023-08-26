package graph

import (
	"context"
	"net/http"
	"net/url"

	"github.com/chiyoi/apricot/kitsune"
)

func Get(ctx context.Context, endpoint string, token string, query ...string) (re *http.Response, err error) {
	u, err := url.JoinPath(endpoint, query...)
	if err != nil {
		return
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return
	}
	kitsune.SetAuthorization(r.Header, token)
	return http.DefaultClient.Do(r)
}
