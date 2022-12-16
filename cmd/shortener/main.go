package main

import (
	"log"
	"net/http"

	"github.com/kerelape/urlshortener/internal/app"
)

const URLShortenerPath = "/"

func main() {
	var shortener = app.NewDatabaseShortener(app.NewFakeDatabase())
	var shortenerHTTPInterface = app.NewMethodFilter(
		http.MethodPost,
		app.NewShortenHandler(shortener),
		app.NewMethodFilter(
			http.MethodGet,
			app.NewRevealHandler(URLShortenerPath, shortener),
			app.MethodNotAllowedHandler(),
		),
	)
	var service = http.NewServeMux()
	service.Handle(URLShortenerPath, shortenerHTTPInterface)
	log.Fatal(http.ListenAndServe(":8080", service))
}
