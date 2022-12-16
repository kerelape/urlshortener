package app

type Shortener interface {
	Shorten(origin string) string
}
