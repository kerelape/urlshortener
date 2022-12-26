package main

import (
	"fmt"
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
	var database = model.NewFakeDatabase()
	var alphabet = model.NewJoinedAlphabet(
		model.NewASCIIAlphabet(48, 57),
		model.NewJoinedAlphabet(
			model.NewASCIIAlphabet(65, 90),
			model.NewASCIIAlphabet(97, 122),
		),
	)
	var shortener = model.NewVerboseShortener(
		model.NewAlphabetShortener(database, alphabet),
		log,
	)
	var router = chi.NewRouter()
	router.Mount(
		"/",
		ui.URLShortener(
			model.NewURLShortener(
				shortener,
				fmt.Sprintf("http://%s%s", Host, Path),
			),
		),
	)
	http.ListenAndServe(Host, router)
}
