package ui

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
)

type API struct {
	shortener model.Shortener
}

type APIRequest struct {
	URL string `json:"url" valid:"url"`
}

type APIResponse struct {
	Result string `json:"result"`
}

func NewAPI(shortener model.Shortener) *API {
	return &API{shortener}
}

func (api *API) Route() http.Handler {
	var router = chi.NewRouter()
	router.Post("/shorten", api.handleShorten)
	return router
}

func (api *API) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}
	var body, readBodyError = io.ReadAll(r.Body)
	if readBodyError != nil {
		http.Error(w, readBodyError.Error(), http.StatusBadRequest)
		return
	}
	var req APIRequest
	var unmarshalError = json.Unmarshal(body, &req)
	if unmarshalError != nil {
		http.Error(w, unmarshalError.Error(), http.StatusBadRequest)
		return
	}
	var shortURL = api.shortener.Shorten(req.URL)
	var resp, marhsalRespError = json.Marshal(APIResponse{Result: shortURL})
	if marhsalRespError != nil {
		http.Error(w, marhsalRespError.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.Itoa(len(resp)))
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
