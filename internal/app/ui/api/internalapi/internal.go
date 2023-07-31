package internalapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/ui/api/internalapi/stats"
)

// Internal is internal api endpoint.
type Internal struct {
	stats stats.Stats
}

// MakeInternal returns a new Internal.
func MakeInternal(statsProvider stats.StatsProvider) Internal {
	return Internal{
		stats: stats.MakeStats(statsProvider),
	}
}

// Route routes Internal.
func (i *Internal) Route() http.Handler {
	router := chi.NewRouter()
	router.Mount("/stats", i.stats.Route())
	return router
}
