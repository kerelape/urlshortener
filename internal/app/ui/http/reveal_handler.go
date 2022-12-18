package http

import (
	"net/http"

	"github.com/kerelape/urlshortener/internal/app/model"
)

type RevealHandler struct {
	Shortener model.Shortener
	Parser    RevealRequestParser
}

// Return new RevealHandler.
func NewRevealHandler(shortener model.Shortener, parser RevealRequestParser) *RevealHandler {
	var handler = new(RevealHandler)
	handler.Shortener = shortener
	handler.Parser = parser
	return handler
}

func (handler *RevealHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var origin, err = handler.Shortener.Reveal(handler.Parser.ParseShortURL(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Add("Location", origin)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
