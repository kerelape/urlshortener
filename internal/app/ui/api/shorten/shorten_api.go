package shorten

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/ui"
)

type ShortenAPI struct {
	shortener model.Shortener
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

func NewShortenAPI(shortener model.Shortener) *ShortenAPI {
	return &ShortenAPI{
		shortener: shortener,
		batch:     NewBatchAPI(shortener),
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
	shortURL, shortenError := shorten.shortener.Shorten(r.Context(), req.URL)
	if shortenError != nil {
		var duplicateError model.DuplicateURLError
		if errors.As(shortenError, &duplicateError) {
			resp, marhsalRespError := json.Marshal(shortenResponse{Result: duplicateError.Origin})
			if marhsalRespError != nil {
				http.Error(w, marhsalRespError.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Content-Length", strconv.Itoa(len(resp)))
			w.WriteHeader(http.StatusConflict)
			w.Write(resp)
			return
		}
		http.Error(w, shortenError.Error(), http.StatusInternalServerError)
		return
	}
	resp, marhsalRespError := json.Marshal(shortenResponse{Result: shortURL})
	if marhsalRespError != nil {
		http.Error(w, marhsalRespError.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.Itoa(len(resp)))
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
