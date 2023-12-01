package authentication

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/chiyoi/apricot/neko"
	"github.com/chiyoi/iter/res"
)

// Login opens an interactive authorization code flow.
// Arguments for LoginURL and RedeemCode are needed.
func Login(ctx context.Context, endpoint Endpoint, listenAddr string, config Config) (token Token, err error) {
	switch runtime.GOOS {
	case "darwin":
		return darwinLogin(ctx, endpoint, listenAddr, config)
	default:
		// TODO: Add device flow for other platforms.
		err = errors.New("unsupported platform")
	}
	return
}

func darwinLogin(ctx context.Context, endpoint Endpoint, listenAddr string, config Config) (token Token, err error) {
	type Values struct {
		token Token
		err   error
	}

	c := make(chan Values)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, err := GetCode(r)
		token, err := res.R(code, err, RedeemCode(endpoint, config))
		c <- Values{token, err}
		if err != nil {
			neko.InternalServerError(w)
		}
		fmt.Fprintln(w, "Login success.")
	})

	srv := &http.Server{
		Addr:    listenAddr,
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
	case v := <-c:
		return v.token, v.err
	}
	return
}
