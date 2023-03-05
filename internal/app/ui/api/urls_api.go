package api

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	body, marshalError := json.Marshal(response)
	if marshalError != nil {
		http.Error(w, marshalError.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

type historyNode struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}
