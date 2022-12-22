package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ChiRevealRequestParser struct {
	ShortURLParamName string
}

func NewChiRevealRequestParser(shortURLParamName string) *ChiRevealRequestParser {
	return &ChiRevealRequestParser{
		ShortURLParamName: shortURLParamName,
	}
}

func (parser *ChiRevealRequestParser) ParseShortURL(request *http.Request) string {
	return chi.URLParam(request, parser.ShortURLParamName)
}
