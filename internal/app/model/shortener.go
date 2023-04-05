package model

import "fmt"

type Shortener interface {
	Shorten(origin string) (string, error)
	Reveal(shortened string) (string, error)
	ShortenAll(origins []string) ([]string, error)
	RevealAll(shortened []string) ([]string, error)
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
