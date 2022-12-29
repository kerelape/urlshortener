package ui

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
)

const ShortURLParam = "short"

type App struct {
	Shortener model.Shortener
}

func NewApp(shortener model.Shortener) *App {
	return &App{Shortener: shortener}
}

func (app *App) Route() http.Handler {
	var router = chi.NewRouter()
	router.Get(fmt.Sprintf("/{%s}", ShortURLParam), app.handleReveal)
	router.Post("/", app.handleShorten)
	return router
}

func (app *App) handleReveal(w http.ResponseWriter, r *http.Request) {
	var short = chi.URLParam(r, ShortURLParam)
	var origin, err = app.Shortener.Reveal(short)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Location", origin)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (app *App) handleShorten(w http.ResponseWriter, r *http.Request) {
	var origin, err = io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var url = string(origin)
	if len(url) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var short = app.Shortener.Shorten(url)
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(short)))
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, short)
}
