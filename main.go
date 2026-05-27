package main

import (
	"log"
	"net"
)

func main() {
	addr := net.UDPAddr{Port: 5454}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("DNS server listening on :%d\n", 5454)
}
