package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/ui"
)

type UserAPI struct {
	urls ui.Entry
}

func NewUserAPI(urls ui.Entry) *UserAPI {
	return &UserAPI{
		urls: urls,
	}
}

func (api *UserAPI) Route() http.Handler {
	router := chi.NewRouter()
	router.Mount("/urls", api.urls.Route())
	return router
}
