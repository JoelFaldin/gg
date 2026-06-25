package server

import (
	"fmt"
	"gg/internal/blocklist"
	"gg/internal/parser"
	"gg/internal/store"
	"log/slog"
	"net"
)

type Server struct {
	conn      *net.UDPConn
	store     *store.Store
	blocklist *blocklist.BlockStruct
}

func StartServer(conn *net.UDPConn, store *store.Store, blocklist *blocklist.BlockStruct) *Server {
	return &Server{conn: conn, store: store, blocklist: blocklist}
}

// Execute the core loop of the server. Read incoming request into a 512-sized slice
func (s *Server) Run() {
	for {
		buf := make([]byte, 4096)

		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		if n < 12 {
			slog.Warn(fmt.Sprintf("Poorly formatted or very short packet received from %s", addr))
			continue
		}

		slog.Debug(fmt.Sprintf("Bytes sent: %d", n))

		// Pass buf[:n], only used bytes
		go s.handle(buf[:n], addr, s.conn, s.store)
	}
}

func (s *Server) handle(buf []byte, addr *net.UDPAddr, conn *net.UDPConn, store *store.Store) {
	res, question, err := parser.ParseMessage(buf, store, conn, addr, s.blocklist)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if len(res) > 0 {
		_, err := conn.WriteToUDP(res, addr)
		if err != nil {
			slog.Error(fmt.Sprintf("Error processing %s: %v", addr, err))
		}

		slog.Info(fmt.Sprintf("Success: %s", question.QName))
	}
}
