package main

import (
	"log"
	"net/http"

	"github.com/kerelape/urlshortener/internal/app"
)

const URLShortenerPath = "/"

func main() {
	var shortener = app.NewDatabaseShortener(app.NewFakeDatabase())
	var shortenerHttpInterface = app.NewMethodFilter(
		http.MethodPost,
		app.NewShortenHandler(shortener),
		app.NewMethodFilter(
			http.MethodGet,
			app.NewRevealHandler(URLShortenerPath, shortener),
			app.MethodNotAllowedHandler(),
		),
	)
	var service = http.NewServeMux()
	service.Handle(URLShortenerPath, shortenerHttpInterface)
	log.Fatal(http.ListenAndServe(":8080", service))
}
