package model

import "strings"

type URLShortener struct {
	Shortener Shortener
	BaseURL   string
	Path      string
}

func NewURLShortener(origin Shortener, baseURL string, path string) *URLShortener {
	return &URLShortener{
		Shortener: origin,
		BaseURL:   baseURL,
		Path:      path,
	}
}

func (shortener *URLShortener) Shorten(origin string) (string, error) {
	var short, shortenError = shortener.Shortener.Shorten(origin)
	return shortener.BaseURL + shortener.Path + short, shortenError
}

func (shortener *URLShortener) Reveal(shortened string) (string, error) {
	return shortener.Shortener.Reveal(strings.TrimPrefix(shortened, shortener.BaseURL+shortener.Path))
}
