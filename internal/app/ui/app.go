package ui

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
)

const ShortURLParam = "short"

type App struct {
	shortener model.Shortener
	history   storage.History
}

func NewApp(shortener model.Shortener, history storage.History) *App {
	return &App{
		shortener: shortener,
		history:   history,
	}
}

func (application *App) Route() http.Handler {
	router := chi.NewRouter()
	router.Get(fmt.Sprintf("/{%s}", ShortURLParam), application.handleReveal)
	router.Post("/", application.handleShorten)
	return router
}

func (application *App) handleReveal(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, ShortURLParam)
	origin, err := application.shortener.Reveal(short)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Location", origin)
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (application *App) handleShorten(w http.ResponseWriter, r *http.Request) {
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
	short, shortenError := application.shortener.Shorten(url)
	if shortenError != nil {
		// duplicate := &model.DuplicateURLError{}
		// if errors.As(shortenError, duplicate) {
		// 	w.WriteHeader(http.StatusConflict)
		// 	io.WriteString(w, duplicate.Origin)
		// 	return
		// }
		http.Error(w, shortenError.Error(), http.StatusInternalServerError)
		return
	}
	user, getTokenError := app.GetToken(r)
	if getTokenError != nil {
		http.Error(w, "No token", http.StatusUnauthorized)
		return
	}
	recordError := application.history.Record(
		user,
		&storage.HistoryNode{
			OriginalURL: url,
			ShortURL:    short,
		},
	)
	if recordError != nil {
		http.Error(w, recordError.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(short)))
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, short)
}
