package token

import (
	"errors"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

// DefaultSecret 测试使用
const (
	DefaultSecret = "jwt-1sNzdiSgnNuxyq2g7xml2JvLArU"
)

// Token Token
type Token struct {
	Account string `json:"acc,omitempty"`
	App     string `json:"app,omitempty"`
	Exp     int64  `json:"exp,omitempty"`
}

var errExpiredToken = errors.New("expired token")

// Valid Valid
func (t *Token) Valid() error {
	if t.Exp < time.Now().Unix() {
		return errExpiredToken
	}
	return nil
}

// Parse ParseJwtToken
func Parse(secret, tk string) (*Token, error) {
	var token = new(Token)
	_, err := jwtgo.ParseWithClaims(tk, token, func(jwttk *jwtgo.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Generate a JWT token
func Generate(secret string, token *Token) (string, error) {
	jtk := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, token)
	return jtk.SignedString([]byte(secret))
}
