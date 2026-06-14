package main

import (
	"fmt"
	"gg/internal/config"
	"gg/internal/logger"
	"gg/internal/server"
	"gg/internal/store"
	"net"

	"github.com/joho/godotenv"
)

func main() {
	// Init logger:
	customLogger := logger.CustomLogger()

	// Load env variables:
	err := godotenv.Load(".env")
	if err != nil {
		customLogger.Error(fmt.Sprintf("Error loading env file: %v", err))
	}

	port := config.GetPort()

	addr := net.UDPAddr{Port: port}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		customLogger.Error(fmt.Sprintf("Error starting app: %v", err))
	}
	defer conn.Close()

	customLogger.Info(fmt.Sprintf("DNS server listening on port :%d", port))

	file := config.GetFile()
	config, err := store.LoadConfig(file)
	if err != nil {
		customLogger.Error(fmt.Sprintf("Could not load data.yaml: %v", err))
	}

	store := store.NewStore(config)

	sv := server.StartServer(conn, store)
	sv.Run()
}
