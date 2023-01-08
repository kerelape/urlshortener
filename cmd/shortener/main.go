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
	DefaultHost = "localhost:8080"
	Path        = "/"
)

const (
	HostEnvironment = "SERVER_ADRESS"
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
	var host = os.Getenv(HostEnvironment)
	if host == "" {
		host = DefaultHost
	}
	var urlShortener = model.NewURLShortener(shortener, "http://"+host+Path)
	var service = chi.NewRouter()
	service.Mount("/", ui.NewApp(urlShortener).Route())
	service.Mount("/api", ui.NewAPI(urlShortener).Route())
	http.ListenAndServe(DefaultHost, service)
}
