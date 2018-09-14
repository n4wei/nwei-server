package api

import (
	"net/http"

	"github.com/n4wei/nwei-server/api/healthcheck"
	"github.com/n4wei/nwei-server/api/weight"
	"github.com/n4wei/nwei-server/db"
	"github.com/n4wei/nwei-server/lib/logger"
)

type controller struct {
	dbClient db.Client
	router   *http.ServeMux
	logger   logger.Logger
}

func NewController(dbClient db.Client, logger logger.Logger) *controller {
	router := http.NewServeMux()
	router.Handle("/healthcheck", chain(healthcheck.Handler, WithLogging(logger)))
	router.Handle("/weight", chain(weight.Handler, WithLogging(logger), WithDB(dbClient)))

	return &controller{
		dbClient: dbClient,
		router:   router,
		logger:   logger,
	}
}

func (c *controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.router.ServeHTTP(w, r)
}
