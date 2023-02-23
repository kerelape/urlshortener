package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/ui"
)

type API struct {
	shorten ui.Entry
	user    ui.Entry
}

func NewAPI(shorten ui.Entry, user ui.Entry) *API {
	return &API{
		shorten: shorten,
		user:    user,
	}
}

func (api *API) Route() http.Handler {
	router := chi.NewRouter()
	router.Mount("/shorten", api.shorten.Route())
	router.Mount("/user", api.user.Route())
	return router
}
