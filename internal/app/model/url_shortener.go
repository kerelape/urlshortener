package model

import (
	"context"
	"errors"
	"strings"

	"github.com/kerelape/urlshortener/internal/app"
)

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

func (shortener *URLShortener) Shorten(ctx context.Context, user app.Token, origin string) (string, error) {
	short, shortenError := shortener.Shortener.Shorten(ctx, user, origin)
	var duplicate DuplicateURLError
	if errors.As(shortenError, &duplicate) {
		duplicate.Origin = shortener.BaseURL + shortener.Path + duplicate.Origin
		return "", duplicate
	}
	return shortener.BaseURL + shortener.Path + short, shortenError
}

func (shortener *URLShortener) Reveal(ctx context.Context, shortened string) (string, error) {
	return shortener.Shortener.Reveal(ctx, strings.TrimPrefix(shortened, shortener.BaseURL+shortener.Path))
}

func (shortener *URLShortener) ShortenAll(ctx context.Context, user app.Token, origins []string) ([]string, error) {
	shorts, shortenError := shortener.Shortener.ShortenAll(ctx, user, origins)
	if shortenError != nil {
		return nil, shortenError
	}
	for i, short := range shorts {
		shorts[i] = shortener.BaseURL + shortener.Path + short
	}
	return shorts, nil
}

func (shortener *URLShortener) RevealAll(ctx context.Context, shortened []string) ([]string, error) {
	for i, short := range shortened {
		shortened[i] = strings.TrimPrefix(short, shortener.BaseURL+shortener.Path)
	}
	origins, revealError := shortener.Shortener.RevealAll(ctx, shortened)
	if revealError != nil {
		return nil, revealError
	}
	return origins, nil
}

func (shortener *URLShortener) Delete(ctx context.Context, user app.Token, shortened []string) error {
	for i, short := range shortened {
		shortened[i] = strings.TrimPrefix(short, shortener.BaseURL+shortener.Path)
	}
	return shortener.Shortener.Delete(ctx, user, shortened)
}
