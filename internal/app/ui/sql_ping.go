package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
)

// SQLPing is the ping end-point.
type SQLPing struct {
	database storage.Database
}

// NewSQLPing returns a new SQLPing.
func NewSQLPing(database storage.Database) *SQLPing {
	return &SQLPing{
		database: database,
	}
}

// Route routes this Entry.
func (ping *SQLPing) Route() http.Handler {
	router := chi.NewRouter()
	router.Get(
		"/",
		func(rw http.ResponseWriter, r *http.Request) {
			status := http.StatusOK
			pingError := ping.database.Ping(r.Context())
			if pingError != nil {
				status = http.StatusInternalServerError
			}
			http.Error(rw, http.StatusText(status), status)
		},
	)
	return router
}
