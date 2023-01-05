package ui

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
)

type Api struct {
	shortener model.Shortener
}

type ApiRequest struct {
	Url string `json:"url" valid:"url"`
}

type ApiResponse struct {
	Result string `json:"result"`
}

func NewApi(shortener model.Shortener) *Api {
	return &Api{shortener}
}

func (api *Api) Route() http.Handler {
	var router = chi.NewRouter()
	router.Post("/shorten", api.handleShorten)
	return router
}

func (api *Api) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}
	var body, readBodyError = io.ReadAll(r.Body)
	if readBodyError != nil {
		http.Error(w, readBodyError.Error(), http.StatusBadRequest)
		return
	}
	var req ApiRequest
	var unmarshalError = json.Unmarshal(body, &req)
	if unmarshalError != nil {
		http.Error(w, unmarshalError.Error(), http.StatusBadRequest)
		return
	}
	var shortURL = api.shortener.Shorten(req.Url)
	var resp, marhsalRespError = json.Marshal(ApiResponse{Result: shortURL})
	if marhsalRespError != nil {
		http.Error(w, marhsalRespError.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.Itoa(len(resp)))
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
