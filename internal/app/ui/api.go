package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/api"
	"github.com/kerelape/urlshortener/internal/app/model"
)

type API struct {
	shortener model.Shortener
	shorten   Entry
}

func NewAPI(shortener model.Shortener) *API {
	return &API{
		shortener: shortener,
		shorten:   api.NewShortenAPI(shortener),
	}
}

func (api *API) Route() http.Handler {
	var router = chi.NewRouter()
	router.Mount("/shorten", api.shorten.Route())
	return router
}
