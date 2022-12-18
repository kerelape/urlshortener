package main

import (
	"log"

	. "github.com/kerelape/urlshortener/internal/app/model"
	. "github.com/kerelape/urlshortener/internal/app/ui/http"
)

func main() {
	log.Fatal(
		NewURLShortenerHTTPService(
			NewDatabaseShortener(
				NewFakeDatabase(),
			),
			"http://localhost:8080",
			"/",
		).Execute(),
	)
}
