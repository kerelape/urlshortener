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
		&Middleware{
			Origin: app.NewShortenerHTTPInterface(
				app.NewDatabaseShortener(app.NewFakeDatabase()),
				URLShortenerPath,
			),
			Middleware: func(w http.ResponseWriter, r *http.Request) {
				println(r.Method, r.Host, r.Proto, r.RequestURI, r.URL)
			},
		},
	)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Middleware struct {
	Origin     http.Handler
	Middleware http.HandlerFunc
}

func (middleware *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware.Middleware(w, r)
	middleware.Origin.ServeHTTP(w, r)
}
