package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/n4wei/nwei-server/lib/logger"
)

const (
	defaultTLSPort = 8443
)

type ServerConfig struct {
	Port    int
	Handler http.Handler
	Logger  logger.Logger

	TLSCertPath  string
	TLSKeyPath   string
	ClientCAPath string
}

type Server struct {
	port    int
	handler http.Handler
	logger  logger.Logger

	tlsCertPath  string
	tlsKeyPath   string
	clientCAPath string
}

func NewServer(config ServerConfig) *Server {
	return &Server{
		port:         config.Port,
		handler:      config.Handler,
		logger:       config.Logger,
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

	if s.port == 0 {
		s.port = defaultTLSPort
	}

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", s.port),
		Handler:   s.handler,
		TLSConfig: tlsConfig,
	}

	s.logger.Printf("listening on %s", server.Addr)

	err = server.ListenAndServeTLS(s.tlsCertPath, s.tlsKeyPath)
	if err != nil {
		return err
	}

	return nil
}
