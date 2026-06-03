package main

import (
	"gg/internal/server"
	"gg/internal/store"
	"log"
	"net"
)

func main() {
	addr := net.UDPAddr{Port: 53}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("DNS server listening on :%d\n", 53)

	config, err := store.LoadConfig("data.yaml")
	if err != nil {
		log.Fatalf("couldnt load data.yaml: %v", err)
	}

	store := store.NewsStore(config)

	sv := server.StartServer(conn, store)
	sv.Run()
}
