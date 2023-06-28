package model

import (
	"context"
	"errors"
	"strings"

	"github.com/kerelape/urlshortener/internal/app"
)

// URLShortener is a shortener that shortens URLs.
type URLShortener struct {
	shortener Shortener
	baseURL   string
	path      string
}

// NewURLShortener returns a new URLShortener.
func NewURLShortener(origin Shortener, baseURL string, path string) *URLShortener {
	return &URLShortener{
		shortener: origin,
		baseURL:   baseURL,
		path:      path,
	}
}

// Shorten shortens the given origin string.
func (shortener *URLShortener) Shorten(ctx context.Context, user app.Token, origin string) (string, error) {
	short, shortenError := shortener.shortener.Shorten(ctx, user, origin)
	var duplicate DuplicateURLError
	if errors.As(shortenError, &duplicate) {
		duplicate.Origin = shortener.baseURL + shortener.path + duplicate.Origin
		return "", duplicate
	}
	return shortener.baseURL + shortener.path + short, shortenError
}

// Reveal returns the original string by the shortened.
func (shortener *URLShortener) Reveal(ctx context.Context, shortened string) (string, error) {
	return shortener.shortener.Reveal(ctx, strings.TrimPrefix(shortened, shortener.baseURL+shortener.path))
}

// ShortenAll shortens a slice of strings and returns
// a slice of short string in the same order.
func (shortener *URLShortener) ShortenAll(ctx context.Context, user app.Token, origins []string) ([]string, error) {
	shorts, shortenError := shortener.shortener.ShortenAll(ctx, user, origins)
	if shortenError != nil {
		return nil, shortenError
	}
	for i, short := range shorts {
		shorts[i] = shortener.baseURL + shortener.path + short
	}
	return shorts, nil
}

// RevealAll returns a slice of original strings in the same order
// as in shortened.
func (shortener *URLShortener) RevealAll(ctx context.Context, shortened []string) ([]string, error) {
	for i, short := range shortened {
		shortened[i] = strings.TrimPrefix(short, shortener.baseURL+shortener.path)
	}
	origins, revealError := shortener.shortener.RevealAll(ctx, shortened)
	if revealError != nil {
		return nil, revealError
	}
	return origins, nil
}

// Delete deletes a string from the shortener.
func (shortener *URLShortener) Delete(ctx context.Context, user app.Token, shortened []string) error {
	for i, short := range shortened {
		shortened[i] = strings.TrimPrefix(short, shortener.baseURL+shortener.path)
	}
	return shortener.shortener.Delete(ctx, user, shortened)
}
