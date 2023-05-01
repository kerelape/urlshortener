package user

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	cjson "github.com/kerelape/cjson/pkg"
	"github.com/kerelape/urlshortener/internal/app"
	"github.com/kerelape/urlshortener/internal/app/model"
)

type URLsAPI struct {
	history model.History
}

func NewURLsAPI(history model.History) URLsAPI {
	return URLsAPI{history}
}

func (h URLsAPI) Route() http.Handler {
	router := chi.NewRouter()
	router.Get("/", h.History)
	router.Delete("/", h.Delete)
	return router
}

func (urls URLsAPI) History(w http.ResponseWriter, r *http.Request) {
	records, recordsError := urls.history.GetRecordsByUser(
		r.Context(),
		r.Context().Value(model.ContextUserTokenKey).(app.Token),
	)
	if recordsError != nil {
		panic(recordsError)
	}
	if len(records) == 0 {
		http.Error(w, "No urls", http.StatusNoContent)
		return
	}
	var response cjson.ArrayBranch = cjson.NewArray()
	for _, record := range records {
		node := cjson.NewObject().
			With("short_url", cjson.NewString(record.To)).
			With("original_url", cjson.NewString(record.From))
		response = response.With(node)
	}
	body, marshalError := response.MarshalJSON()
	if marshalError != nil {
		http.Error(w, marshalError.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (urls URLsAPI) Delete(w http.ResponseWriter, r *http.Request) {
	token, tokenErr := app.GetToken(r)
	if tokenErr != nil {
		status := http.StatusUnauthorized
		http.Error(w, http.StatusText(status), status)
		return
	}
	body, bodyError := io.ReadAll(r.Body)
	if bodyError != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	request, requestError := cjson.NewDocument(body).Parse()
	if requestError != nil {
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}
	if !request.Array().Present() {
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}
	requestedForRemoval := request.Array().Value()
	records, recordsError := urls.history.GetRecordsByUser(r.Context(), token)
	if recordsError != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	var own cjson.ObjectBranch = cjson.NewObject()
	for _, record := range records {
		own = own.With(record.To, cjson.NewString(record.From))
	}
	for i := 0; i < requestedForRemoval.Size(); i++ {
		node := requestedForRemoval.At(i).Value()
		if !node.String().Present() {
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
			return
		}
		removee := node.String().Value().Content()
		if !own.Found(removee).Present() {
			status := http.StatusUnauthorized
			http.Error(w, http.StatusText(status), status)
			return
		}

	}
	w.WriteHeader(http.StatusAccepted)
}
