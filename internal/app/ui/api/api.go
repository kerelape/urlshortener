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

// API is api end-point.
type API struct {
	shorten ui.Entry
	user    ui.Entry
}

// NewAPI returns a new API.
func NewAPI(shortener model.Shortener, history storage.History) *API {
	return &API{
		shorten: shorten.NewShortenAPI(shortener, history),
		user:    user.NewUserAPI(history, shortener),
	}
}

// Route routes this Entry.
func (api *API) Route() http.Handler {
	router := chi.NewRouter()
	router.Mount("/shorten", api.shorten.Route())
	router.Mount("/user", api.user.Route())
	return router
}
