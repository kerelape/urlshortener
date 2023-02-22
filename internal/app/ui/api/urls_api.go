package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type URLSAPI struct {
}

func NewURLSAPI() *URLSAPI {
	return &URLSAPI{}
}

func (api *URLSAPI) Route() http.Handler {
	var router = chi.NewRouter()
	return router
}
