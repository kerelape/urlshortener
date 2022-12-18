package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	. "github.com/kerelape/urlshortener/internal/app/model"
)

type URLShortenerHTTPService struct {
	Shortener Shortener
	Host      string
	Path      string
}

func NewURLShortenerHTTPService(shortener Shortener, host string, path string) *URLShortenerHTTPService {
	var service = new(URLShortenerHTTPService)
	service.Shortener = shortener
	service.Host = host
	service.Path = path
	return service
}

func (service *URLShortenerHTTPService) Execute() error {
	var router = chi.NewRouter()
	router.Route(service.Path, func(router chi.Router) {
		var shortener = NewURLShortener(
			service.Shortener,
			service.Host+service.Path,
		)
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
