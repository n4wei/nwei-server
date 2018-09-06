package controller

import (
	"net/http"

	"github.com/n4wei/nwei-server/lib/logger"
)

type Adapter func(http.Handler) http.Handler

func Adapt(handler http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		handler = adapter(handler)
	}
	return handler
}

func WithLogging(logger logger.Logger) Adapter {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%#v\n", r)
			handler.ServeHTTP(w, r)
		})
	}
}
