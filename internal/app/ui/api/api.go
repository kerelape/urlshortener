package api

import (
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
	"github.com/kerelape/urlshortener/internal/app/ui"
	"github.com/kerelape/urlshortener/internal/app/ui/api/internalapi"
	"github.com/kerelape/urlshortener/internal/app/ui/api/internalapi/stats"
	"github.com/kerelape/urlshortener/internal/app/ui/api/shorten"
	"github.com/kerelape/urlshortener/internal/app/ui/api/user"
)

// API is api end-point.
type API struct {
	shorten  ui.Entry
	user     ui.Entry
	internal internalapi.Internal
}

// NewAPI returns a new API.
func NewAPI(
	shortener model.Shortener,
	history storage.History,
	statsProvider stats.StatsProvider,
	internalTrustedSubnet *net.IPNet,
) *API {
	return &API{
		shorten:  shorten.NewShortenAPI(shortener, history),
		user:     user.NewUserAPI(history, shortener),
		internal: internalapi.MakeInternal(statsProvider, internalTrustedSubnet),
	}
}

// Route routes this Entry.
func (api *API) Route() http.Handler {
	router := chi.NewRouter()
	router.Mount("/shorten", api.shorten.Route())
	router.Mount("/user", api.user.Route())
	router.Mount("/internal", api.internal.Route())
	return router
}
