package graph

import (
	"context"
	"net/http"
	"net/url"

	"github.com/chiyoi/apricot/kitsune"
)

func Read(ctx context.Context, endpoint string, token string, v any, query ...string) (err error) {
	u, err := url.JoinPath(endpoint, query...)
	if err != nil {
		return
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return
	}
	kitsune.SetAuthorization(r.Header, token)

	re, err := http.DefaultClient.Do(r)
	if err != nil {
		return
	}
	defer re.Body.Close()

	return kitsune.ParseResponse(re, v)
}
