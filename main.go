package main

import (
	"gg/internal/config"
	"gg/internal/server"
	"gg/internal/store"
	"log"
	"net"

	"github.com/joho/godotenv"
)

func main() {
	// Load env variables:
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading env file:", err)
	}

	port := config.GetPort()

	addr := net.UDPAddr{Port: port}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("DNS server listening on :%d\n", port)

	file := config.GetFile()
	config, err := store.LoadConfig(file)
	if err != nil {
		log.Fatalf("couldnt load data.yaml: %v", err)
	}

	store := store.NewStore(config)

	sv := server.StartServer(conn, store)
	sv.Run()
}
