package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
)

type URLShortenerHTTPService struct {
	Shortener model.Shortener
	Host      string
	Path      string
	Log       model.Log
}

func NewURLShortenerHTTPService(shortener model.Shortener, host string, path string) *URLShortenerHTTPService {
	return &URLShortenerHTTPService{
		Shortener: shortener,
		Host:      host,
		Path:      path,
	}
}

func (service *URLShortenerHTTPService) Execute() error {
	var router = chi.NewRouter()
	var shortener = model.NewURLShortener(
		service.Shortener,
		fmt.Sprintf("http://%s%s", service.Host, service.Path),
	)
	router.Route(service.Path, func(router chi.Router) {
		router.Get(
			"/{short}",
			handlerToHandlerFunc(
				NewRevealHandler(
					shortener,
					NewChiRevealRequestParser("short"),
				),
			),
		)
		router.Post(
			"/",
			handlerToHandlerFunc(
				NewShortenHandler(shortener),
			),
		)
	})
	return http.ListenAndServe(service.Host, router)
}

func handlerToHandlerFunc(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
