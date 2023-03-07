package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
)

type SQLPing struct {
	database storage.Database
}

func NewSQLPing(database storage.Database) *SQLPing {
	return &SQLPing{
		database: database,
	}
}

func (ping *SQLPing) Route() http.Handler {
	router := chi.NewRouter()
	router.Get(
		"/",
		func(rw http.ResponseWriter, r *http.Request) {
			pingError := ping.database.Ping()
			if pingError != nil {
				status := http.StatusInternalServerError
				http.Error(rw, http.StatusText(status), status)
			} else {
				status := http.StatusOK
				http.Error(rw, http.StatusText(status), status)
			}
		},
	)
	return router
}
