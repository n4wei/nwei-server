package controller

import (
	"net"

	"github.com/n4wei/nwei-server/lib/logger"
)

func HandleConn(conn net.Conn, logger logger.Logger) {
	defer conn.Close()
}
