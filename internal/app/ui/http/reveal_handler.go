package http

import (
	"net/http"

	. "github.com/kerelape/urlshortener/internal/app/model"
)

type RevealHandler struct {
	Shortener Shortener
}

// Return new RevealHandler.
func NewRevealHandler(shortener Shortener) *RevealHandler {
	var handler = new(RevealHandler)
	handler.Shortener = shortener
	return handler
}

func (handler *RevealHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Get real short url
	var origin, err = handler.Shortener.Reveal("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Location", origin)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
