package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/ui"
)

type API struct {
	shorten ui.Entry
}

func NewAPI(shorten ui.Entry) *API {
	return &API{
		shorten: shorten,
	}
}

func (api *API) Route() http.Handler {
	var router = chi.NewRouter()
	router.Mount("/shorten", api.shorten.Route())
	return router
}
