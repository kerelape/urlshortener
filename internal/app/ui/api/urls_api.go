package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserURLs struct {
}

func NewUserURLs() *UserURLs {
	return &UserURLs{}
}

func (api *UserURLs) Route() http.Handler {
	var router = chi.NewRouter()
	return router
}
