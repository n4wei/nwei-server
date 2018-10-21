package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ServerConfig struct {
	Port    int
	Handler http.Handler

	TLSCertPath  string
	TLSKeyPath   string
	ClientCAPath string
}

func NewServer(config ServerConfig) (*http.Server, error) {
	clientCA, err := ioutil.ReadFile(config.ClientCAPath)
	if err != nil {
		return nil, err
	}

	clientCAPool := x509.NewCertPool()
	if ok := clientCAPool.AppendCertsFromPEM(clientCA); !ok {
		return nil, errors.New("failed to add client CA to pool")
	}

	tlsConfig := &tls.Config{
		ClientCAs:  clientCAPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// TODO: timeouts
	return &http.Server{
		Addr:      fmt.Sprintf(":%d", config.Port),
		Handler:   config.Handler,
		TLSConfig: tlsConfig,
	}, nil
}
