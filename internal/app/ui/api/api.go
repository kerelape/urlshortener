package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/ui"
)

type API struct {
	shortener model.Shortener
	shorten   ui.Entry
}

func NewAPI(shortener model.Shortener) *API {
	return &API{
		shortener: shortener,
		shorten:   NewShortenAPI(shortener),
	}
}

func (api *API) Route() http.Handler {
	var router = chi.NewRouter()
	router.Mount("/shorten", api.shorten.Route())
	return router
}
