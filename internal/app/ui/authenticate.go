package ui

import (
	"context"
	"net/http"

	"github.com/kerelape/urlshortener/internal/app"
	"github.com/kerelape/urlshortener/internal/app/model"
)

func Authenticate() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, getTokenError := app.GetToken(r)
			if getTokenError != nil {
				status := http.StatusUnauthorized
				http.Error(w, http.StatusText(status), status)
				return
			}
			next.ServeHTTP(
				w,
				r.WithContext(
					context.WithValue(
						r.Context(),
						model.ContextUserTokenKey,
						user,
					),
				),
			)
		})
	}
}
