package server

import (
	"crypto/tls"

	"github.com/n4wei/nwei-server/controller"
	"github.com/n4wei/nwei-server/lib/logger"
)

const (
	defaultTLSPort = ":8443"
)

type ServerConfig struct {
	TLSCertPath string
	TLSKeyPath  string
	Logger      logger.Logger
}

type Server struct {
	tlsCertPath string
	tlsKeyPath  string
	Logger      logger.Logger
}

func NewServer(config ServerConfig) *Server {
	return &Server{
		tlsCertPath: config.TLSCertPath,
		tlsKeyPath:  config.TLSKeyPath,
		Logger:      config.Logger,
	}
}

func (s *Server) Serve() error {
	cert, err := tls.LoadX509KeyPair(s.tlsCertPath, s.tlsKeyPath)
	if err != nil {
		return err
	}

	listener, err := tls.Listen("tcp", defaultTLSPort, &tls.Config{Certificates: []tls.Certificate{cert}})
	if err != nil {
		return err
	}
	defer listener.Close()

	s.Logger.Printf("nwei-server listening on %s", defaultTLSPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.Logger.Error(err)
			continue
		}
		go controller.HandleConn(conn, s.Logger)
	}
}
