package http

import (
	"io"
	"net/http"

	"github.com/kerelape/urlshortener/internal/app/model"
)

type ShortenHandler struct {
	Shortener model.Shortener
}

// Return new ShortenHandler.
func NewShortenHandler(shortener model.Shortener) *ShortenHandler {
	return &ShortenHandler{
		Shortener: shortener,
	}
}

func (handler *ShortenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read url", http.StatusInternalServerError)
		return
	}
	var url = string(body)
	if url == "" {
		http.Error(w, "No URL provided", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, handler.Shortener.Shorten(url))
}
