package config

import (
	"fmt"
	"gg/internal/logger"
	"log"
	"os"
	"strconv"
)

func GetPort() int {
	p := os.Getenv("PORT")
	port, err := strconv.Atoi(p)
	if err != nil {
		customLogger := logger.CustomLogger()
		customLogger.Error(fmt.Sprintf("Port conversion failed: %v", err))
	}

	if port == 0 {
		port = 8053
	}

	return port
}

func GetFile() string {
	file := os.Getenv("FILE")

	if file == "" {
		log.Fatal("Invalid file location, check .env file.")
	}

	return file
}

func GetAddress() string {
	address := os.Getenv("ADDRESS")

	if address == "" {
		log.Fatal("Invalid address value. Remember to specify the port (eg. 8.8.8.8:53)")
	}

	return address
}
