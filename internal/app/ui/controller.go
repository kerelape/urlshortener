package ui

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
)

func URLShortenerApp(shortener model.Shortener, host string, path string) http.Handler {
	var app = chi.NewRouter()
	var urlShortener = model.NewURLShortener(shortener, fmt.Sprintf("http://%s%s", host, path))
	app.Mount(path, URLShortener(urlShortener))
	return app
}

func URLShortener(shortener model.Shortener) http.Handler {
	var router = chi.NewRouter()
	router.Get("/{short}", revealHandler(shortener))
	router.Post("/", shortenHandler(shortener))
	return router
}

func revealHandler(shortener model.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var short = chi.URLParam(r, "short")
		var origin, err = shortener.Reveal(short)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Add("Location", origin)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func shortenHandler(shortener model.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var origin, err = io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var url = string(origin)
		if len(url) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var short = shortener.Shorten(url)
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(short)))
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, short)
	}
}
