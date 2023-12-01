package graph

import (
	"context"
	"io"
	"net/url"

	"github.com/chiyoi/apricot/kitsune"
	"github.com/chiyoi/iter/res"
)

const (
	Endpoint = "https://graph.microsoft.com/"
)

func Read(ctx context.Context, token string, item any, query ...string) (err error) {
	u, err := url.JoinPath(Endpoint, query...)
	return res.C(u, err, kitsune.GetJSON(ctx, item, kitsune.SetAuthorizationHook(token)))
}

func Request(ctx context.Context, token string, query ...string) (body io.ReadCloser, err error) {
	u, err := url.JoinPath(Endpoint, query...)
	return res.R(u, err, kitsune.GetStream(ctx, kitsune.SetAuthorizationHook(token)))
}
