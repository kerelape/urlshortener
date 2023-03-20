package model

type Shortener interface {
	Shorten(origin string) (string, error)
	Reveal(shortened string) (string, error)
	ShortenAll(origins []string) ([]string, error)
	RevealAll(shortened []string) ([]string, error)
}
