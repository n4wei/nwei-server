package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/n4wei/nwei-server/api"
	"github.com/n4wei/nwei-server/db"
	"github.com/n4wei/nwei-server/lib/logger"
	"github.com/n4wei/nwei-server/server"
)

const (
	cleanupAndShutdownTimeout = 5 * time.Second
)

func main() {
	var serverConfig server.ServerConfig
	flag.IntVar(&serverConfig.Port, "port", 0, "The port that the server will listen on")
	flag.StringVar(&serverConfig.TLSCertPath, "tls-cert", "", "The filepath to the certificate used for TLS")
	flag.StringVar(&serverConfig.TLSKeyPath, "tls-key", "", "The filepath to the private key used for TLS")
	flag.StringVar(&serverConfig.ClientCAPath, "client-ca", "", "The filepath to the client's CA certificate")

	var dbConfig db.DBConfig
	flag.StringVar(&dbConfig.URL, "db-url", "localhost:27017", "The full database URL with optional auth")

	flag.Parse()

	logger := logger.NewLogger()

	dbConfig.Logger = logger
	dbClient, err := db.NewClient(dbConfig)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	serverConfig.Handler = api.NewController(dbClient, logger).Handler()
	server, err := server.NewServer(serverConfig)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-stop
		logger.Printf("caught signal: %v", sig)
		logger.Print("shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), cleanupAndShutdownTimeout)
		defer cancel()

		err := dbClient.Close(ctx)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

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
