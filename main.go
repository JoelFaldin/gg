package main

import (
	"flag"
	"fmt"
	"gg/internal/config"
	"gg/internal/logger"
	"gg/internal/server"
	"gg/internal/store"
	"log/slog"
	"net"

	"github.com/joho/godotenv"
)

func main() {
	// Flags:
	verbose := flag.Bool("verbose", false, "enable verbose (debug) logging")
	flag.Parse()

	// Init logger:
	logLevel := slog.LevelInfo
	if *verbose {
		logLevel = slog.LevelDebug
	}
	customLogger := logger.CustomLogger(logLevel)
	slog.SetDefault(customLogger)

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
