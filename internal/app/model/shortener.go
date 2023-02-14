package model

type Shortener interface {
	Shorten(origin string) (string, error)
	Reveal(shortened string) (string, error)
}
