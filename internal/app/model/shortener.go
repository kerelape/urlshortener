package model

import (
	"context"
	"fmt"

	"github.com/kerelape/urlshortener/internal/app"
)

type Shortener interface {
	Shorten(ctx context.Context, user app.Token, origin string) (string, error)
	Reveal(ctx context.Context, shortened string) (string, error)
	ShortenAll(ctx context.Context, user app.Token, origins []string) ([]string, error)
	RevealAll(ctx context.Context, shortened []string) ([]string, error)
	Delete(ctx context.Context, user app.Token, shortened []string) error
}

type DuplicateURLError struct {
	Origin string
}

func NewDuplicateURLError(origin string) DuplicateURLError {
	return DuplicateURLError{
		Origin: origin,
	}
}

func (e DuplicateURLError) Error() string {
	return fmt.Sprintf("duplicate URL: %s", e.Origin)
}
