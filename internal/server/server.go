package server

import (
	"fmt"
	"gg/internal/parser"
	"gg/internal/store"
	"log"
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
	for {
		buf := make([]byte, 512)

		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		if n < 12 {
			fmt.Printf("paquete mal formateado o muy corto recibido desde %s\n", addr)
			continue
		}

		fmt.Printf("bytes enviados: %d\n", n)

		// Pass buf[:n], only used bytes
		go s.handle(buf[:n], addr, s.conn, s.store)
	}
}

func (s *Server) handle(buf []byte, addr *net.UDPAddr, conn *net.UDPConn, store *store.Store) {
	_, err := parser.ParseMessage(buf, store)
	if err != nil {
		log.Println("parse error:", err)
	}

	conn.WriteToUDP([]byte("ok"), addr)
}
