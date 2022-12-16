package app

import "net/http"

func NewShortenerHTTPInterface(shortener Shortener, prefix string) http.Handler {
	return NewMethodFilter(
		http.MethodPost,
		NewShortenHandler(shortener),
		NewMethodFilter(
			http.MethodGet,
			NewRevealHandler(prefix, shortener),
			MethodNotAllowedHandler(),
		),
	)
}
