package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// WebUI is a general purpose one-level composite routes.
type WebUI struct {
	routes map[string]Entry
}

// NewWebUI returns a new WebUI.
func NewWebUI(routes map[string]Entry) *WebUI {
	return &WebUI{
		routes: routes,
	}
}

// Route routes this Entry.
func (ui *WebUI) Route() http.Handler {
	router := chi.NewRouter()
	for path, entry := range ui.routes {
		router.Mount(path, entry.Route())
	}
	return router
}
