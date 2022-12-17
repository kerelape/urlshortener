package main

import (
	"log"
	"net/http"

	"github.com/kerelape/urlshortener/internal/app"
)

const URLShortenerPath = "/"

func main() {
	http.Handle(
		URLShortenerPath,
		app.NewShortenerHTTPInterface(
			app.NewDatabaseShortener(app.NewFakeDatabase()),
			URLShortenerPath,
			"http://localhost:8080",
		),
	)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
