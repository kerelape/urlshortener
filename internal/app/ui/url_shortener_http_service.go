package ui

import (
	"fmt"
	"io"
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
		router.Get("/{short}", func(w http.ResponseWriter, r *http.Request) {
			var short = chi.URLParam(r, "short")
			var origin, err = shortener.Reveal(short)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Add("Location", origin)
			w.WriteHeader(http.StatusTemporaryRedirect)
		})
		router.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var origin, err = io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var url = string(origin)
			if len(url) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var short = shortener.Shorten(url)
			w.Header().Add("Content-Length", fmt.Sprintf("%d", len(short)))
			w.WriteHeader(http.StatusCreated)
			io.WriteString(w, short)
		})
	})
	return http.ListenAndServe(service.Host, router)
}
