package ui

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
)

const shortURLParam = "short"

// App is the root REST endpoint.
type App struct {
	shortener model.Shortener
	history   storage.History
}

// NewApp returns a new App.
func NewApp(shortener model.Shortener, history storage.History) *App {
	return &App{
		shortener: shortener,
		history:   history,
	}
}

// Route routes this Entry.
func (application *App) Route() http.Handler {
	router := chi.NewRouter()
	router.Get(fmt.Sprintf("/{%s}", shortURLParam), application.HandleReveal)
	router.Post("/", application.HandleShorten)
	return router
}

// HandleReveal redirects to the original URL behind the short url name.
func (application *App) HandleReveal(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, shortURLParam)
	origin, err := application.shortener.Reveal(r.Context(), short)
	if err != nil {
		if errors.Is(err, storage.ErrValueDeleted) {
			w.WriteHeader(http.StatusGone)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Location", origin)
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// HandleShorten shortens a URL provided in the request body.
func (application *App) HandleShorten(w http.ResponseWriter, r *http.Request) {
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
	user, getTokenError := app.GetToken(r)
	if getTokenError != nil {
		http.Error(w, "No token", http.StatusUnauthorized)
		return
	}
	short, shortenError := application.shortener.Shorten(r.Context(), user, url)
	if shortenError != nil {
		duplicate := &model.DuplicateURLError{}
		if errors.As(shortenError, duplicate) {
			w.WriteHeader(http.StatusConflict)
			io.WriteString(w, duplicate.Origin)
			return
		}
		http.Error(w, shortenError.Error(), http.StatusInternalServerError)
		return
	}
	recordError := application.history.Record(
		r.Context(),
		user,
		storage.HistoryNode{
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
