package main

import (
	"log"

	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/ui/http"
)

func main() {
	log.Fatal(
		http.NewURLShortenerHTTPService(
			model.NewAlphabetShortener(
				model.NewFakeDatabase(),
				model.NewJoinedAlphabet(
					model.NewASCIIAlphabet(48, 57),
					model.NewJoinedAlphabet(
						model.NewASCIIAlphabet(65, 90),
						model.NewASCIIAlphabet(97, 122),
					),
				),
			),
			"localhost:8080",
			"/",
		).Execute(),
	)
}
