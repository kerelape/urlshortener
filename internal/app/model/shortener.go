package model

import (
	"context"
	"fmt"

	"github.com/kerelape/urlshortener/internal/app"
)

type Shortener interface {
	// Shorten shortens the given origin string.
	Shorten(ctx context.Context, user app.Token, origin string) (string, error)

	// Reveal returns the original string by the shortened.
	Reveal(ctx context.Context, shortened string) (string, error)

	// ShortenAll shortens a slice of strings and returns
	// a slice of short string in the same order.
	ShortenAll(ctx context.Context, user app.Token, origins []string) ([]string, error)

	// RevealAll returns a slice of original strings in the same order
	// as in shortened.
	RevealAll(ctx context.Context, shortened []string) ([]string, error)

	// Delete deletes a string from the shortener.
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
