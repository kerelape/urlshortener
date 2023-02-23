package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserURLs struct{}

func NewUserURLs() *UserURLs {
	return &UserURLs{}
}

func (api *UserURLs) Route() http.Handler {
	router := chi.NewRouter()
	router.Get("/", api.ServeHTTP)
	return router
}

func (api *UserURLs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
