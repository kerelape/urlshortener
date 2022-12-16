package app

import "net/http"

func UrlShortenerHandler(path string, shortener Shortener) http.Handler {
	return NewMethodFilter(
		http.MethodPost,
		NewShortenHandler(shortener),
		NewMethodFilter(
			http.MethodGet,
			NewRevealHandler(path, shortener),
			MethodNotAllowedHandler(),
		),
	)
}
