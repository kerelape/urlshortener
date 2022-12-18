package http

import "net/http"

type RevealRequestParser interface {
	ParseShortURL(request *http.Request) string
}
