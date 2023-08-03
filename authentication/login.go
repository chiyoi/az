package authentication

import (
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
func Login(endpoint Endpoint, config Config) (token Token, err error) {
	u, err := url.Parse(config.RedirectURI)
	if err != nil {
		return
	}

	t, e := make(chan Token), make(chan error)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", u.Port()),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code, err := GetCode(r)
			if err != nil {
				neko.InternalServerError(w, "Login failed.")
				e <- err
				return
			}

			token, err := RedeemCode(code, endpoint, config)
			if err != nil {
				neko.InternalServerError(w, "Get token failed.")
				e <- err
				return
			}

			fmt.Fprintln(w, "Login success.")
			t <- token
		}),
	}

	go neko.StartServer(srv, false)
	defer neko.StopServer(srv)

	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", LoginURL(endpoint, config)).Start()
	default:
		err = errors.New("unsupported platform")
	}
	if err != nil {
		return
	}

	select {
	case token = <-t:
	case err = <-e:
	}
	return
}
