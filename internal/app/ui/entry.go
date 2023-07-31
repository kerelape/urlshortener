package ui

import "net/http"

// Entry is a API end-point.
type Entry interface {
	// Route routes this Entry and returns an http.Handler.
	Route() http.Handler
}
