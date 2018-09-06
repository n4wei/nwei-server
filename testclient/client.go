package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	defaultServerName = "nwei-server"
)

func main() {
	var clientConfig ClientConfig
	flag.StringVar(&clientConfig.URL, "url", "", "The URL of the server")
	flag.StringVar(&clientConfig.TLSCertPath, "tls-cert", "", "The filepath to the certificate used for TLS")
	flag.StringVar(&clientConfig.TLSKeyPath, "tls-key", "", "The filepath to the private key used for TLS")
	flag.StringVar(&clientConfig.ServerCAPath, "server-ca", "", "The filepath to the server's CA certificate")
	flag.Parse()

	c, err := NewClient(clientConfig)
	if err != nil {
		HandleErr(err)
	}

	resp, err := c.httpClient.Get(clientConfig.URL)
	if err != nil {
		HandleErr(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		HandleErr(err)
	}
	defer resp.Body.Close()

	fmt.Printf("%s", body)
}

func HandleErr(err error) {
	fmt.Fprint(os.Stderr, err.Error())
	os.Exit(1)
}

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
