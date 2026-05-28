package server

import (
	"gg/internal/parser"
	"log"
	"net"
)

type Server struct {
	conn *net.UDPConn
}

func StartServer(conn *net.UDPConn) *Server {
	return &Server{conn: conn}
}

// Execute the core loop of the server. Read incoming request into a 512-sized slice
func (s *Server) Run() {
	for {
		buf := make([]byte, 512)

		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		// Pass buf[:n], only used bytes
		go s.handle(buf[:n], addr, s.conn)
	}
}

func (s *Server) handle(buf []byte, addr *net.UDPAddr, conn *net.UDPConn) {
	_, err := parser.ParseMessage(buf)
	if err != nil {
		log.Println("parse error:", err)
	}

	conn.WriteToUDP([]byte("ok"), addr)
}
