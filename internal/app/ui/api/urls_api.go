package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
)

type UserURLs struct {
	history storage.History
}

func NewUserURLs(history storage.History) *UserURLs {
	return &UserURLs{
		history: history,
	}
}

func (api *UserURLs) Route() http.Handler {
	router := chi.NewRouter()
	router.Get("/", api.ServeHTTP)
	return router
}

func (api *UserURLs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tokenCookie, tokenCookieError := r.Cookie("token")
	if tokenCookieError != nil {
		panic(tokenCookieError)
	}
	token, tokenError := app.TokenFromString(tokenCookie.Value)
	if tokenError != nil {
		panic(tokenError)
	}
	records, recordsError := api.history.GetRecordsByUser(token)
	if recordsError != nil {
		panic(recordsError)
	}
	if len(records) == 0 {
		http.Error(w, "No urls", http.StatusNoContent)
		return
	}
	response := make([]*historyNode, len(records))
	for i := 0; i < len(response); i++ {
		response[i] = &historyNode{
			Short:    records[i].ShortURL,
			Original: records[i].OriginalURL,
		}
	}
	writer := json.NewEncoder(w)
	encodeError := writer.Encode(response)
	if encodeError != nil {
		http.Error(w, encodeError.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

type historyNode struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}
