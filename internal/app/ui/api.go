package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/api"
	"github.com/kerelape/urlshortener/internal/app/model"
)

const (
	shortenPath string = "/shorten"
)

type API struct {
	shortener model.Shortener
	shorten   http.Handler
}

func NewAPI(shortener model.Shortener) *API {
	return &API{
		shortener: shortener,
		shorten:   api.NewShortenAPI(shortener),
	}
}

func (api *API) Route() http.Handler {
	var router = chi.NewRouter()
	router.Post(shortenPath, api.shorten.ServeHTTP)
	return router
}
