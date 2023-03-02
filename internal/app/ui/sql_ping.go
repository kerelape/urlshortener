package ui

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SQLPing struct {
	database *sql.DB
}

func NewSQLPing(database *sql.DB) *SQLPing {
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
