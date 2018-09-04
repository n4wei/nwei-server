package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/n4wei/nwei-server/controller"
	"github.com/n4wei/nwei-server/lib/logger"
)

const (
	defaultTLSPort = ":8443"
)

type ServerConfig struct {
	Logger logger.Logger

	TLSCertPath  string
	TLSKeyPath   string
	ClientCAPath string
}

type Server struct {
	Logger logger.Logger

	tlsCertPath  string
	tlsKeyPath   string
	clientCAPath string
}

func NewServer(config ServerConfig) *Server {
	return &Server{
		Logger:       config.Logger,
		tlsCertPath:  config.TLSCertPath,
		tlsKeyPath:   config.TLSKeyPath,
		clientCAPath: config.ClientCAPath,
	}
}

func (s *Server) Serve() error {
	clientCA, err := ioutil.ReadFile(s.clientCAPath)
	if err != nil {
		return err
	}

	clientCAPool := x509.NewCertPool()
	if ok := clientCAPool.AppendCertsFromPEM(clientCA); !ok {
		return errors.New("failed to add client CA to pool")
	}

	tlsConfig := &tls.Config{
		ClientCAs:  clientCAPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      defaultTLSPort,
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/", controller.Handler)
	s.Logger.Printf("listening on %s", defaultTLSPort)

	err = server.ListenAndServeTLS(s.tlsCertPath, s.tlsKeyPath)
	if err != nil {
		return err
	}

	return nil
}
