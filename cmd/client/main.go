package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/n4wei/nwei-server/lib/logger"
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
