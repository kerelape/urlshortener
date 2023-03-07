package shorten

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
	"github.com/kerelape/urlshortener/internal/app/ui"
)

type ShortenAPI struct {
	shortener model.Shortener
	history   storage.History
	batch     ui.Entry
}

type (
	shortenRequest struct {
		URL string `json:"url" valid:"url"`
	}

	shortenResponse struct {
		Result string `json:"result"`
	}
)

func NewShortenAPI(shortener model.Shortener, history storage.History) *ShortenAPI {
	return &ShortenAPI{
		shortener: shortener,
		history:   history,
		batch:     NewBatchAPI(shortener, history),
	}
}

func (shorten *ShortenAPI) Route() http.Handler {
	router := chi.NewRouter()
	router.Post("/", shorten.ServeHTTP)
	router.Mount("/batch", shorten.batch.Route())
	return router
}

func (shorten *ShortenAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}
	body, readBodyError := io.ReadAll(r.Body)
	if readBodyError != nil {
		http.Error(w, readBodyError.Error(), http.StatusBadRequest)
		return
	}
	var req shortenRequest
	unmarshalError := json.Unmarshal(body, &req)
	if unmarshalError != nil {
		http.Error(w, unmarshalError.Error(), http.StatusBadRequest)
		return
	}
	shortURL, shortenError := shorten.shortener.Shorten(req.URL)
	if shortenError != nil {
		http.Error(w, shortenError.Error(), http.StatusInternalServerError)
		return
	}
	resp, marhsalRespError := json.Marshal(shortenResponse{Result: shortURL})
	if marhsalRespError != nil {
		http.Error(w, marhsalRespError.Error(), http.StatusInternalServerError)
		return
	}
	user, getTokenError := app.GetToken(r)
	if getTokenError != nil {
		http.Error(w, "No token", http.StatusUnauthorized)
		return
	}
	recordError := shorten.history.Record(
		user,
		&storage.HistoryNode{
			OriginalURL: req.URL,
			ShortURL:    shortURL,
		},
	)
	if recordError != nil {
		http.Error(w, recordError.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.Itoa(len(resp)))
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
