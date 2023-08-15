package internalapi

import (
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kerelape/urlshortener/internal/app/ui/api/internalapi/stats"
)

// Internal is internal api endpoint.
type Internal struct {
	stats stats.Stats

	trustedSubnet *net.IPNet
}

// MakeInternal returns a new Internal.
func MakeInternal(statsProvider stats.StatsProvider, trustedSubnet *net.IPNet) Internal {
	return Internal{
		stats:         stats.MakeStats(statsProvider),
		trustedSubnet: trustedSubnet,
	}
}

// Route routes Internal.
func (i *Internal) Route() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(out http.ResponseWriter, in *http.Request) {
			var realIP, _, ipErr = net.ParseCIDR(in.RemoteAddr)
			if ipErr != nil {
				var status = http.StatusBadRequest
				http.Error(out, http.StatusText(status), status)
				return
			}
			if !i.trustedSubnet.Contains(realIP) {
				var status = http.StatusBadRequest
				http.Error(out, http.StatusText(status), status)
				return
			}
			next.ServeHTTP(out, in)
		})
	})
	router.Mount("/stats", i.stats.Route())
	return router
}
