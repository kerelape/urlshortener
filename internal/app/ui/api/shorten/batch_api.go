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

func (api *BatchAPI) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	user, userError := app.GetToken(r)
	if userError != nil {
		status := http.StatusUnauthorized
		http.Error(rw, http.StatusText(status), status)
		return
	}
	var request []*batchAPIRequestNode
	decodeError := json.NewDecoder(r.Body).Decode(&request)
	if decodeError != nil {
		status := http.StatusBadRequest
		http.Error(rw, http.StatusText(status), status)
		return
	}
	response := make([]*batchAPIResponseNode, len(request))
	for i := 0; i < len(response); i++ {
		requestNode := request[i]
		shortURL, shortenError := api.shortener.Shorten(requestNode.OriginalURL)
		if shortenError != nil {
			status := http.StatusInternalServerError
			http.Error(rw, http.StatusText(status), status)
			return
		}
		response[i] = &batchAPIResponseNode{
			CorrelationID: requestNode.CorrelationID,
			ShortURL:      shortURL,
		}
	}
	for i := 0; i < len(response); i++ {
		recordError := api.history.Record(user, &storage.HistoryNode{
			OriginalURL: request[i].OriginalURL,
			ShortURL:    response[i].ShortURL,
		})
		if recordError != nil {
			http.Error(rw, recordError.Error(), http.StatusInternalServerError)
			return
		}
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
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
