package model

import "strings"

type URLShortener struct {
	Shortener Shortener
	BaseURL   string
}

func NewURLShortener(origin Shortener, baseURL string) *URLShortener {
	return &URLShortener{
		Shortener: origin,
		BaseURL:   baseURL,
	}
}

func (shortener *URLShortener) Shorten(origin string) string {
	return shortener.BaseURL + shortener.Shortener.Shorten(origin)
}

func (shortener *URLShortener) Reveal(shortened string) (string, error) {
	return shortener.Shortener.Reveal(strings.TrimPrefix(shortened, shortener.BaseURL))
}
