package controller

import (
	"net/http"

	"github.com/n4wei/nwei-server/lib/logger"
)

type controller struct {
	router *http.ServeMux
	logger logger.Logger
}

func NewController(logger logger.Logger) *controller {
	return &controller{
		router: http.NewServeMux(),
		logger: logger,
	}
}

func (c *controller) Handler() http.Handler {
	c.router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success\n"))
	})
	return c.router
}
