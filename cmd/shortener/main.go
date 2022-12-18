package main

import (
	"log"
	"net/http"

	. "github.com/kerelape/urlshortener/internal/app/model"
	. "github.com/kerelape/urlshortener/internal/app/ui/http"
)

const URLShortenerPath = "/"

func main() {
	http.Handle(
		URLShortenerPath,
		NewShortenerHTTPInterface(
			NewDatabaseShortener(NewFakeDatabase()),
			URLShortenerPath,
			"http://localhost:8080",
		),
	)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
