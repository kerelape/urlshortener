package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/ui"
)

const (
	Host = "localhost:8080"
	Path = "/"
)

func main() {
	var log = model.NewFormattedLog(
		model.NewWriterLog(os.Stdout, os.Stderr),
		time.UnixDate,
	)
	var shortener = model.NewVerboseShortener(
		model.NewAlphabetShortener(
			model.NewFakeDatabase(),
			model.NewBase62Alphabet(),
		),
		log,
	)
	var urlShortener = model.NewURLShortener(shortener, "http://"+Host+Path)
	var service = chi.NewRouter()
	service.Mount("/", ui.NewApp(urlShortener).Route())
	service.Mount("/api", ui.NewApi(urlShortener).Route())
	http.ListenAndServe(Host, service)
}
