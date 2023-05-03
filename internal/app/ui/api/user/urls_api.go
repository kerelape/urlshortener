package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
)

type URLsAPI struct {
	history   storage.History
	shortener model.Shortener
}

func NewURLsAPI(history storage.History, shortener model.Shortener) *URLsAPI {
	return &URLsAPI{
		history:   history,
		shortener: shortener,
	}
}

func (api *URLsAPI) Route() http.Handler {
	router := chi.NewRouter()
	router.Get("/", api.History)
	router.Delete("/", api.Delete)
	return router
}

func (api *URLsAPI) History(w http.ResponseWriter, r *http.Request) {
	tokenCookie, tokenCookieError := r.Cookie("token")
	if tokenCookieError != nil {
		panic(tokenCookieError)
	}
	token, tokenError := app.TokenFromString(tokenCookie.Value)
	if tokenError != nil {
		panic(tokenError)
	}
	records, recordsError := api.history.GetRecordsByUser(r.Context(), token)
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
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (api *URLsAPI) Delete(w http.ResponseWriter, r *http.Request) {
	user, userError := app.GetToken(r)
	if userError != nil {
		status := http.StatusUnauthorized
		http.Error(w, http.StatusText(status), status)
		return
	}
	urls := []string{}
	urlsError := json.NewDecoder(r.Body).Decode(&urls)
	if urlsError != nil {
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}
	go api.shortener.Delete(context.Background(), user, urls)
	w.WriteHeader(http.StatusAccepted)
}

type historyNode struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}
