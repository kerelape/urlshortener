package main

import (
	"os"
	"time"

	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/ui/http"
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
	var service = http.NewURLShortenerHTTPService(
		shortener,
		"localhost:8080",
		"/",
		log,
	)
	service.Execute()
}
