package model

import (
	"context"
	"fmt"
)

type Shortener interface {
	Shorten(ctx context.Context, origin string) (string, error)
	Reveal(ctx context.Context, shortened string) (string, error)
	ShortenAll(ctx context.Context, origins []string) ([]string, error)
	RevealAll(ctx context.Context, shortened []string) ([]string, error)
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
