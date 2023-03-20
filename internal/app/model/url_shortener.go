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
	short, shortenError := shortener.Shortener.Shorten(origin)
	return shortener.BaseURL + shortener.Path + short, shortenError
}

func (shortener *URLShortener) Reveal(shortened string) (string, error) {
	return shortener.Shortener.Reveal(strings.TrimPrefix(shortened, shortener.BaseURL+shortener.Path))
}

func (shortener *URLShortener) ShortenAll(origins []string) ([]string, error) {
	shorts, shortenError := shortener.Shortener.ShortenAll(origins)
	if shortenError != nil {
		return nil, shortenError
	}
	for i, short := range shorts {
		shorts[i] = shortener.BaseURL + shortener.Path + short
	}
	return shorts, nil
}

func (shortener *URLShortener) RevealAll(shortened []string) ([]string, error) {
	for i, short := range shortened {
		shortened[i] = strings.TrimPrefix(short, shortener.BaseURL+shortener.Path)
	}
	origins, revealError := shortener.Shortener.RevealAll(shortened)
	if revealError != nil {
		return nil, revealError
	}
	return origins, nil
}
