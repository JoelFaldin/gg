package main

import (
	"gg/internal/server"
	"log"
	"net"
)

func main() {
	addr := net.UDPAddr{Port: 8053}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("DNS server listening on :%d\n", 8053)

	sv := server.StartServer(conn)
	sv.Run()
}
