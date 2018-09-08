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
	"strings"

	"github.com/n4wei/nwei-server/lib/logger"
)

const (
	defaultServerName = "nwei-server"
)

func main() {
	var postData string
	var clientConfig ClientConfig
	flag.StringVar(&clientConfig.URL, "url", "", "The URL of the server")
	flag.StringVar(&clientConfig.TLSCertPath, "tls-cert", "", "The filepath to the certificate used for TLS")
	flag.StringVar(&clientConfig.TLSKeyPath, "tls-key", "", "The filepath to the private key used for TLS")
	flag.StringVar(&clientConfig.ServerCAPath, "server-ca", "", "The filepath to the server's CA certificate")
	flag.StringVar(&postData, "post-data", "", "Data sent on a post request")
	flag.Parse()

	c, err := NewClient(clientConfig)
	if err != nil {
		HandleErr(err)
	}

	var resp *http.Response
	if postData != "" {
		resp, err = c.httpClient.Post(clientConfig.URL, "application/json", strings.NewReader(postData))
		if err != nil {
			HandleErr(err)
		}
	} else {
		resp, err = c.httpClient.Get(clientConfig.URL)
		if err != nil {
			HandleErr(err)
		}
	}

	logger := logger.NewLogger()
	logger.Print(logger.FormatHTTPResponse(resp))
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
