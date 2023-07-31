package ui

import (
	"encoding/hex"
	"net/http"

	"github.com/kerelape/urlshortener/internal/app"
)

const tokenCookieName = "token"

// Tokenize is a middleware that assigns tokens to each new user.
func Tokenize() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, cookieError := r.Cookie(tokenCookieName)
			if cookieError != nil {
				token := app.NewToken()
				app.SetToken(w, token)
				r.AddCookie(&http.Cookie{
					Name:  tokenCookieName,
					Value: hex.EncodeToString(token[:]),
				})
			}
			next.ServeHTTP(w, r)
		})
	}
}
