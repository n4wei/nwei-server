package main

import (
	"flag"
	"os"

	"github.com/n4wei/nwei-server/api"
	"github.com/n4wei/nwei-server/db"
	"github.com/n4wei/nwei-server/lib/logger"
	"github.com/n4wei/nwei-server/server"
)

func main() {
	var serverConfig server.ServerConfig
	flag.StringVar(&serverConfig.TLSCertPath, "tls-cert", "", "The filepath to the certificate used for TLS")
	flag.StringVar(&serverConfig.TLSKeyPath, "tls-key", "", "The filepath to the private key used for TLS")
	flag.StringVar(&serverConfig.ClientCAPath, "client-ca", "", "The filepath to the client's CA certificate")

	var dbConfig db.DBConfig
	flag.StringVar(&dbConfig.URL, "db-url", "localhost", "The database URL")
	flag.IntVar(&dbConfig.Port, "db-port", 27017, "The database port")

	flag.Parse()

	logger := logger.NewLogger()

	dbConfig.Logger = logger
	dbClient, err := db.NewClient(dbConfig)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	serverConfig.Logger = logger
	serverConfig.Handler = api.NewController(dbClient, logger).Handler()
	server := server.NewServer(serverConfig)

	// TODO: graceful stop
	err = server.Serve()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
