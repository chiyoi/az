package authentication

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/chiyoi/apricot/neko"
)

// Login opens an interactive authorization code flow.
// Arguments for LoginURL and RedeemCode are needed.
func Login(ctx context.Context, endpoint Endpoint, config Config) (token Token, err error) {
	switch runtime.GOOS {
	case "darwin":
		return darwinLogin(ctx, endpoint, config)
	default:
		// TODO: Add device flow for other platforms.
		err = errors.New("unsupported platform")
	}
	return
}

func darwinLogin(ctx context.Context, endpoint Endpoint, config Config) (token Token, err error) {
	u, err := url.Parse(config.RedirectURI)
	if err != nil {
		return
	}
	t, e := make(chan Token), make(chan error)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, err := GetCode(r)
		if err != nil {
			e <- err
			neko.BadRequest(w)
			return
		}

		token, err := RedeemCode(code, endpoint, config)
		if err != nil {
			e <- err
			neko.InternalServerError(w)
			return
		}

		t <- token
		fmt.Fprintln(w, "Login success.")
	})
	srv := &http.Server{
		Addr:    u.Host,
		Handler: h,
	}

	go neko.StartServer(srv, false)
	defer neko.StopServer(srv)

	if err = exec.Command("open", LoginURL(endpoint, config)).Start(); err != nil {
		return
	}

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case token = <-t:
	case err = <-e:
	}
	return
}
