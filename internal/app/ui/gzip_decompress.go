package ui

import (
	"compress/gzip"
	"io"
	"net/http"
)

type (
	readCloser struct {
		io.Reader
		io.Closer
	}
)

func Decompress() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			request := *r
			if r.Header.Get("Content-Encoding") == "gzip" {
				decompressedBody, decompressedBodyError := gzip.NewReader(r.Body)
				if decompressedBodyError != nil {
					http.Error(w, decompressedBodyError.Error(), http.StatusInternalServerError)
					return
				}
				request.Body = &readCloser{Reader: decompressedBody, Closer: r.Body}
			}
			next.ServeHTTP(w, &request)
		})
	}
}
