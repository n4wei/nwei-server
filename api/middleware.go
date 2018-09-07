package api

import (
	"net/http"

	"github.com/n4wei/nwei-server/lib/logger"
)

type middleware func(http.HandlerFunc) http.HandlerFunc

func chain(handler http.HandlerFunc, middleware ...middleware) http.HandlerFunc {
	for _, m := range middleware {
		handler = m(handler)
	}
	return handler
}

func WithLogging(logger logger.Logger) middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger.Print(logger.FormatHTTPRequest(r))
			handler.ServeHTTP(w, r)
		}
	}
}
