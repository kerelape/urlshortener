package app

import (
	"net/http"
	"strings"
)

type RevealHandler struct {
	Prefix    string
	Shortener Shortener
}

// Return new RevealHandler.
func NewRevealHandler(prefix string, shortener Shortener) *RevealHandler {
	var handler = new(RevealHandler)
	handler.Prefix = prefix
	handler.Shortener = shortener
	return handler
}

func (self *RevealHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var shortened = strings.TrimPrefix(self.Prefix, r.URL.Path)
	var origin, err = self.Shortener.Reveal(shortened)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Location", origin)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
