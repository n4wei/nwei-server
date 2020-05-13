package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/n4wei/nwei-server/api"
	"github.com/n4wei/nwei-server/db/mongo"
	"github.com/n4wei/nwei-server/lib/logger"
)

const (
	defaultCleanupAndShutdownTimeout = 5 * time.Second
)

func main() {
	var serverConfig ServerConfig
	flag.IntVar(&serverConfig.Port, "port", 8443, "port that the server will accept TLS connections on")
	flag.StringVar(&serverConfig.TLSCertPath, "tls-cert", "", "filepath to server TLS certificate")
	flag.StringVar(&serverConfig.TLSKeyPath, "tls-key", "", "filepath to server TLS private key")
	flag.StringVar(&serverConfig.ClientCAPath, "client-ca", "", "filepath to clients' CA certificate for MTLS")

	var dbConfig mongo.DBConfig
	flag.StringVar(&dbConfig.URL, "db-url", "mongodb://localhost:27017", "database URL")

	flag.Parse()

	logger := logger.NewLogger()

	dbConfig.Logger = logger
	dbClient, err := mongo.NewClient(dbConfig)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	serverConfig.Handler = api.NewController(dbClient, logger)
	server, err := NewServer(serverConfig)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-stop
		logger.Printf("caught signal: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), defaultCleanupAndShutdownTimeout)
		defer cancel()

		logger.Print("shutting down server...")
		err = server.Shutdown(ctx)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		os.Exit(0)
	}()

	logger.Printf("listening on %s", server.Addr)

	err = server.ListenAndServeTLS(serverConfig.TLSCertPath, serverConfig.TLSKeyPath)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
