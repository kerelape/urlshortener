package app

import "net/http"

func NewShortenerHTTPInterface(shortener Shortener, prefix string, baseURL string) http.Handler {
	return NewMethodFilter(
		http.MethodPost,
		NewShortenHandler(shortener, baseURL+prefix),
		NewMethodFilter(
			http.MethodGet,
			NewRevealHandler(prefix, shortener),
			MethodNotAllowedHandler(),
		),
	)
}
