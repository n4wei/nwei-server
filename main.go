package main

import (
	"flag"
	"os"

	"github.com/n4wei/nwei-server/controller"
	"github.com/n4wei/nwei-server/lib/logger"
	"github.com/n4wei/nwei-server/server"
)

func main() {
	var serverConfig server.ServerConfig
	flag.StringVar(&serverConfig.TLSCertPath, "tls-cert", "", "The filepath to the certificate used for TLS")
	flag.StringVar(&serverConfig.TLSKeyPath, "tls-key", "", "The filepath to the private key used for TLS")
	flag.StringVar(&serverConfig.ClientCAPath, "client-ca", "", "The filepath to the client's CA certificate")
	flag.Parse()

	logger := logger.NewLogger()
	serverConfig.Handler = controller.NewController(logger).Handler()
	serverConfig.Logger = logger

	server := server.NewServer(serverConfig)

	// TODO: graceful stop
	err := server.Serve()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
