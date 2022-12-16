package main

import (
	"log"
	"net/http"

	"github.com/kerelape/urlshortener/internal/app"
)

const URLShortenerPath = "/"

func main() {
	var service = http.NewServeMux()
	service.Handle(
		URLShortenerPath,
		app.NewShortenerHTTPInterface(
			app.NewDatabaseShortener(app.NewFakeDatabase()),
			URLShortenerPath,
		),
	)
	log.Fatal(http.ListenAndServe(":8080", service))
}
