package app

import (
	"encoding/hex"
	"math/rand"
	"net/http"
	"time"
)

// Token is a user token.
type Token [8]byte

// NewToken returns a new random Token.
func NewToken() Token {
	var token Token
	rand.Read(token[:])
	return token
}

// TokenFromString decodes a hex string into a user token.
func TokenFromString(origin string) (Token, error) {
	bytes, decodeError := hex.DecodeString(origin)
	var token Token
	copy(token[:], bytes)
	return token, decodeError
}

// SetToken sets token cookie.
func SetToken(rw http.ResponseWriter, token Token) {
	http.SetCookie(
		rw,
		&http.Cookie{
			Name:     "token",
			Value:    string(hex.EncodeToString(token[:])),
			Expires:  time.Now().Add(time.Hour * 24 * 60),
			Path:     "/",
			HttpOnly: true,
		},
	)
}

// GetToken returns token from cookies.
func GetToken(r *http.Request) (Token, error) {
	cookie, cookieError := r.Cookie("token")
	if cookieError != nil {
		return *new(Token), cookieError
	}
	token, decodeError := TokenFromString(cookie.Value)
	if decodeError != nil {
		return *new(Token), nil
	}
	return token, nil
}
