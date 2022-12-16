package app

type Shortener interface {
	Shorten(origin string) string
	Reveal(shortened string) (string, error)
}
