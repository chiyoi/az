package authentication

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// LoginURL needs endpoint Endpoint.Authorization, and
// configs Config.ClientID, Config.Scopes and Config.RedirectURI
func LoginURL(endpoint Endpoint, config Config) string {
	q := url.Values{}
	q.Set("client_id", config.ClientID)
	q.Set("response_type", "code")
	q.Set("redirect_uri", config.RedirectURI)
	q.Set("scope", strings.Join(config.Scopes, " "))
	q.Set("prompt", "select_account")
	return fmt.Sprintf("%s?%s", endpoint.Authorize, q.Encode())
}

func GetCode(r *http.Request) (code string, err error) {
	if err := r.URL.Query().Get("error"); err != "" {
		return "", fmt.Errorf("[%s] %s", err, r.URL.Query().Get("error_description"))
	}

	code = r.URL.Query().Get("code")
	if code == "" {
		return "", errors.New("invalid code")
	}
	return
}

// LogoutURL needs endpoint Endpoint.Logout
func LogoutURL(endpoint Endpoint, postLogoutRedirectURI string) string {
	q := url.Values{}
	q.Set("post_logout_redirect_uri", postLogoutRedirectURI)
	return fmt.Sprintf("%s?%s", endpoint.Logout, q.Encode())
}

// RedeemCode needs endpoint Endpoint.Token, and
// configs Config.ClientID, Config.Scopes and RedirectURI
// config.RedirectURI must be the same as acquiring the code
func RedeemCode(code string, endpoint Endpoint, config Config) (tokens Token, err error) {
	q := url.Values{}
	q.Set("client_id", config.ClientID)
	q.Set("scope", strings.Join(config.Scopes, " "))
	q.Set("redirect_uri", config.RedirectURI)
	q.Set("grant_type", "authorization_code")
	q.Set("code", code)

	resp, err := http.DefaultClient.Post(endpoint.Token, "application/x-www-form-urlencoded", strings.NewReader(q.Encode()))
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		var ate acquireTokenError
		if err = json.NewDecoder(resp.Body).Decode(&ate); err != nil {
			return
		}

		err = fmt.Errorf("(%s) %s", ate.Error, ate.ErrorDescription)
		return
	}

	var atr acquireTokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&atr); err != nil {
		return
	}

	return Token{
		AccessToken: atr.AccessToken,
		ExpiresAt:   time.Now().UTC().Add(time.Second * time.Duration(atr.ExpiresIn)).Unix(),

		RefreshToken: atr.RefreshToken,
	}, nil
}

// Refresh needs endpoint Endpoint.Token, and
// configs Config.ClientID and Config.Scopes
func Refresh(refreshToken string, endpoint Endpoint, config Config) (tokens Token, err error) {
	q := url.Values{}
	q.Set("client_id", config.ClientID)
	q.Set("scope", strings.Join(config.Scopes, " "))
	q.Set("grant_type", "refresh_token")
	q.Set("refresh_token", refreshToken)

	resp, err := http.DefaultClient.Post(endpoint.Token, "application/x-www-form-urlencoded", strings.NewReader(q.Encode()))
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		var ate acquireTokenError
		if err = json.NewDecoder(resp.Body).Decode(&ate); err != nil {
			return
		}

		err = fmt.Errorf("(%s) %s", ate.Error, ate.ErrorDescription)
		return
	}

	var atr acquireTokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&atr); err != nil {
		return
	}

	return Token{
		AccessToken: atr.AccessToken,
		ExpiresAt:   time.Now().UTC().Add(time.Second * time.Duration(atr.ExpiresIn)).Unix(),

		RefreshToken: atr.RefreshToken,
	}, nil
}

type Endpoint struct {
	Authorize string
	Token     string
	Logout    string
}

type Config struct {
	ClientID    string
	Scopes      []string
	RedirectURI string
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`

	RefreshToken string `json:"refresh_token"`
}

func (token Token) Expired() bool {
	return time.Now().Unix() > token.ExpiresAt
}

type acquireTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type acquireTokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCodes       []int  `json:"error_codes"`
}
