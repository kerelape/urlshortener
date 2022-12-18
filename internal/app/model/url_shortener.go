package model

import "strings"

type URLShortener struct {
	Shortener Shortener
	Host      string
}

func NewURLShortener(origin Shortener, host string) *URLShortener {
	var shortener = new(URLShortener)
	shortener.Shortener = origin
	shortener.Host = host
	return shortener
}

func (shortener *URLShortener) Shorten(origin string) string {
	return shortener.Host + shortener.Shortener.Shorten(origin)
}

func (shortener *URLShortener) Reveal(shortened string) (string, error) {
	return shortener.Shortener.Reveal(strings.TrimPrefix(shortened, shortener.Host))
}
