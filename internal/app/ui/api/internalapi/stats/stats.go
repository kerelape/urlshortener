package stats

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
)

// Stats is stats api endpoint.
type Stats struct {
	provider StatsProvider
}

// MakeStats returns a new Stats.
func MakeStats(provider StatsProvider) Stats {
	return Stats{
		provider: provider,
	}
}

// Route routes Stats.
func (s *Stats) Route() http.Handler {
	router := chi.NewRouter()
	router.Get("/", s.ServeHTTP)
	return router
}

// ServeHTTP handle requests to Stats.
func (s *Stats) ServeHTTP(out http.ResponseWriter, in *http.Request) {
	group, gtx := errgroup.WithContext(in.Context())
	group.SetLimit(2)
	var urls, users = new(int), new(int)
	group.Go(func() error {
		u, err := s.provider.URLs(gtx)
		*urls = u
		return err
	})
	group.Go(func() error {
		u, err := s.provider.Users(gtx)
		*users = u
		return err
	})
	if err := group.Wait(); err != nil {
		var status = http.StatusInternalServerError
		http.Error(out, http.StatusText(status), status)
		return
	}

	var response = struct {
		URLs  int `json:"urls"`
		Users int `json:"users"`
	}{*urls, *users}
	body, bodyError := json.Marshal(&response)
	if bodyError != nil {
		var status = http.StatusInternalServerError
		http.Error(out, http.StatusText(status), status)
		return
	}

	out.WriteHeader(http.StatusOK)
	out.Write(body)
}
