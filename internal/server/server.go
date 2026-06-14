package server

import (
	"fmt"
	"gg/internal/logger"
	"gg/internal/parser"
	"gg/internal/store"
	"net"
)

type Server struct {
	conn  *net.UDPConn
	store *store.Store
}

func StartServer(conn *net.UDPConn, store *store.Store) *Server {
	return &Server{conn: conn, store: store}
}

// Execute the core loop of the server. Read incoming request into a 512-sized slice
func (s *Server) Run() {
	customLogger := logger.CustomLogger()

	for {
		buf := make([]byte, 4096)

		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		if n < 12 {
			customLogger.Warn(fmt.Sprintf("Poorly formatted or very short packet received from %s", addr))
			continue
		}

		customLogger.Info(fmt.Sprintf("Bytes sent: %d", n))

		// Pass buf[:n], only used bytes
		go s.handle(buf[:n], addr, s.conn, s.store)
	}
}

func (s *Server) handle(buf []byte, addr *net.UDPAddr, conn *net.UDPConn, store *store.Store) {
	customLogger := logger.CustomLogger()

	res, err := parser.ParseMessage(buf, store, conn, addr)
	if err != nil {
		customLogger.Error(err.Error())
		return
	}

	if len(res) > 0 {
		_, err := conn.WriteToUDP(res, addr)
		if err != nil {
			customLogger.Error(fmt.Sprintf("Error processing %s: %v", addr, err))
		}
	}
}
