package controller

import (
	"net/http"

	"github.com/n4wei/nwei-server/controller/health"
	"github.com/n4wei/nwei-server/controller/test"
	"github.com/n4wei/nwei-server/lib/logger"
)

type controller struct {
	router *http.ServeMux
	logger logger.Logger
}

func NewController(logger logger.Logger) *controller {
	router := http.NewServeMux()
	router.Handle("/test", Adapt(test.Handler(), WithLogging(logger)))
	router.Handle("/health", Adapt(health.Handler(), WithLogging(logger)))

	return &controller{
		router: router,
		logger: logger,
	}
}

func (c *controller) Handler() http.Handler {
	return c.router
}
