package main

import (
	"crypto/tls"
	"crypto/x509"
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
	flag.StringVar(&clientConfig.TLSCAPath, "tls-ca", "", "The filepath to the CA certificate for TLS")
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

	fmt.Printf("response: %s\n", body)
}

func HandleErr(err error) {
	fmt.Fprint(os.Stderr, err.Error())
	os.Exit(1)
}

type ClientConfig struct {
	URL       string
	TLSCAPath string
}

type Client struct {
	httpClient *http.Client
}

func NewClient(config ClientConfig) (*Client, error) {
	caPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(config.TLSCAPath)
	if err != nil {
		return nil, err
	}
	if ok := caPool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("error adding CA to pool")
	}

	return &Client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					ServerName: defaultServerName,
					RootCAs:    caPool,
				},
			},
		},
	}, nil
}
