package ui

import (
	"encoding/hex"
	"net/http"

	"github.com/kerelape/urlshortener/internal/app"
)

const tokenCookieName = "token"

func Tokenize() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, cookieError := r.Cookie(tokenCookieName)
			if cookieError != nil {
				token := app.NewToken()
				r.AddCookie(&http.Cookie{
					Name:  tokenCookieName,
					Value: hex.EncodeToString(token[:]),
				})
			}
		})
	}
}
