package app

import (
	"encoding/hex"
	"math/rand"
	"net/http"
)

type Token [8]byte

func NewToken() Token {
	var token Token
	rand.Read(token[:])
	return token
}

func TokenFromString(origin string) (Token, error) {
	bytes, decodeError := hex.DecodeString(origin)
	var token Token
	copy(token[:], bytes)
	return token, decodeError
}

func SetToken(rw http.ResponseWriter, token Token) {
	http.SetCookie(
		rw,
		&http.Cookie{
			Name:  "token",
			Value: string(hex.EncodeToString(token[:])),
		},
	)
}

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
