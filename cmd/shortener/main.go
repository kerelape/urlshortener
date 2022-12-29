package main

import (
	"net/http"
	"os"
	"time"

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
	var app = ui.NewApp(
		model.NewURLShortener(
			shortener,
			"http://"+Host+Path,
		),
	)
	http.ListenAndServe(Host, app.Route())
}
