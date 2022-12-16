package app

import (
	"io"
	"net/http"
)

type ShortenHandler struct {
	Shortener Shortener
}

// Return new ShortenHandler.
func NewShortenHandler(shortener Shortener) *ShortenHandler {
	var handler = new(ShortenHandler)
	handler.Shortener = shortener
	return handler
}

func (self *ShortenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	io.WriteString(w, self.Shortener.Shorten(url))
	w.WriteHeader(http.StatusCreated)
}
