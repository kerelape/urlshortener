package ui

import "net/http"

type Entry interface {
	Route() http.Handler
}
