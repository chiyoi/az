package encrypt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func Assertion(audience string, issuer, subject string, headers map[string]any, key any) (assertion string, err error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Audience:  audience,
		ExpiresAt: now.Add(time.Minute * 10).Unix(),
		Id:        uuid.NewString(),
		IssuedAt:  now.Unix(),
		Issuer:    issuer,
		NotBefore: now.Unix(),
		Subject:   subject,
	})
	for k, v := range headers {
		token.Header[k] = v
	}
	return token.SignedString(key)
}
