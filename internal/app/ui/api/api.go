package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
	"github.com/kerelape/urlshortener/internal/app/ui"
	"github.com/kerelape/urlshortener/internal/app/ui/api/shorten"
	"github.com/kerelape/urlshortener/internal/app/ui/api/user"
)

type API struct {
	shorten ui.Entry
	user    ui.Entry
}

func NewAPI(shortener model.Shortener, history storage.History) *API {
	return &API{
		shorten: shorten.NewShortenAPI(shortener, history),
		user:    user.NewUserAPI(history, shortener),
	}
}

func (api *API) Route() http.Handler {
	router := chi.NewRouter()
	router.Mount("/shorten", api.shorten.Route())
	router.Mount("/user", api.user.Route())
	return router
}
