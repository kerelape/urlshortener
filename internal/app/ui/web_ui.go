package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type WebUI struct {
	routes map[string]Entry
}

func NewWebUI(routes map[string]Entry) *WebUI {
	return &WebUI{
		routes: routes,
	}
}

func (ui *WebUI) Route() http.Handler {
	router := chi.NewRouter()
	for path, entry := range ui.routes {
		router.Mount(path, entry.Route())
	}
	return router
}
