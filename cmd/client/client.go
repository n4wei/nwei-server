package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	defaultServerName = "nwei-server"
)

type ClientConfig struct {
	URL string

	TLSCertPath  string
	TLSKeyPath   string
	ServerCAPath string
}

type Client struct {
	httpClient *http.Client
}

func NewClient(config ClientConfig) (*Client, error) {
	cert, err := tls.LoadX509KeyPair(config.TLSCertPath, config.TLSKeyPath)
	if err != nil {
		return nil, err
	}

	serverCA, err := ioutil.ReadFile(config.ServerCAPath)
	if err != nil {
		return nil, err
	}

	serverCAPool := x509.NewCertPool()
	if ok := serverCAPool.AppendCertsFromPEM(serverCA); !ok {
		return nil, errors.New("failed to add server CA to pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      serverCAPool,
		ServerName:   defaultServerName,
	}
	tlsConfig.BuildNameToCertificate()

	return &Client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
	}, nil
}
