package http

import "net/http"

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
}

func MethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(MethodNotAllowed)
}
