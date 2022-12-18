package http

import (
	"net/http"

	. "github.com/kerelape/urlshortener/internal/app/model"
)

func NewShortenerHTTPInterface(shortener Shortener) http.Handler {
	return NewMethodFilter(
		http.MethodPost,
		NewShortenHandler(shortener),
		NewMethodFilter(
			http.MethodGet,
			NewRevealHandler(shortener),
			MethodNotAllowedHandler(),
		),
	)
}
