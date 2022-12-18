package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	. "github.com/kerelape/urlshortener/internal/app/model"
	. "github.com/kerelape/urlshortener/internal/app/ui/http"
)

func main() {
	var shortener = NewDatabaseShortener(NewFakeDatabase())
	var router = chi.NewRouter()
	router.Route("/", func(router chi.Router) {
		var shortener = NewURLShortener(shortener, "http://localhost:8080/")
		router.Get(
			"/{short}",
			HandlerToHandlerFunc(
				NewRevealHandler(
					shortener,
					NewChiRevealRequestParser("short"),
				),
			),
		)
		router.Post(
			"/",
			HandlerToHandlerFunc(NewShortenHandler(shortener)),
		)
	})
	log.Fatal(http.ListenAndServe(":8080", router))
}

func HandlerToHandlerFunc(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
