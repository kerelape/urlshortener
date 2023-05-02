package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
	"github.com/kerelape/urlshortener/internal/app/ui"
)

type UserAPI struct {
	urls ui.Entry
}

func NewUserAPI(history storage.History, shortener model.Shortener) *UserAPI {
	return &UserAPI{
		urls: NewURLsAPI(history, shortener),
	}
}

func (api *UserAPI) Route() http.Handler {
	router := chi.NewRouter()
	router.Mount("/urls", api.urls.Route())
	return router
}
