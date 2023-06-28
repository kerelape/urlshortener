package shorten

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
)

type BatchAPI struct {
	shortener model.Shortener
	history   storage.History
}

func NewBatchAPI(shortener model.Shortener, history storage.History) *BatchAPI {
	return &BatchAPI{
		shortener: shortener,
		history:   history,
	}
}

func (api *BatchAPI) Route() http.Handler {
	var router = chi.NewRouter()
	router.Post("/", api.ServeHTTP)
	return router
}

// ServeHTTP shortens many URL and returns an array of short ones.
func (api *BatchAPI) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	user, userError := app.GetToken(r)
	if userError != nil {
		status := http.StatusUnauthorized
		http.Error(rw, http.StatusText(status), status)
		return
	}
	var request []batchAPIRequestNode
	decodeError := json.NewDecoder(r.Body).Decode(&request)
	if decodeError != nil {
		status := http.StatusBadRequest
		http.Error(rw, http.StatusText(status), status)
		return
	}
	response := make([]batchAPIResponseNode, len(request))
	origins := make([]string, len(request))
	for i, requestNode := range request {
		origins[i] = requestNode.OriginalURL
	}
	shorts, shortenError := api.shortener.ShortenAll(r.Context(), user, origins)
	if shortenError != nil {
		status := http.StatusInternalServerError
		http.Error(rw, http.StatusText(status), status)
		return
	}
	for i := range response {
		response[i] = batchAPIResponseNode{
			CorrelationID: request[i].CorrelationID,
			ShortURL:      shorts[i],
		}
	}
	for i := 0; i < len(response); i++ {
		recordError := api.history.Record(
			r.Context(),
			user, storage.HistoryNode{
				OriginalURL: request[i].OriginalURL,
				ShortURL:    response[i].ShortURL,
			},
		)
		if recordError != nil {
			http.Error(rw, recordError.Error(), http.StatusInternalServerError)
			return
		}
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(&response)
}

type (
	batchAPIRequestNode struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}
	batchAPIResponseNode struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)
