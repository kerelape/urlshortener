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
	router := chi.NewRouter()
	router.Get(fmt.Sprintf("/{%s}", ShortURLParam), app.handleReveal)
	router.Post("/", app.handleShorten)
	return router
}

func (app *App) handleReveal(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, ShortURLParam)
	origin, err := app.Shortener.Reveal(short)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Location", origin)
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (app *App) handleShorten(w http.ResponseWriter, r *http.Request) {
	origin, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	url := string(origin)
	if len(url) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	short, shortenError := app.Shortener.Shorten(url)
	if shortenError != nil {
		http.Error(w, shortenError.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(short)))
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, short)
}
