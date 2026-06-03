package config

import (
	"log"
	"os"
	"strconv"
)

func GetPort() int {
	p := os.Getenv("PORT")
	port, err := strconv.Atoi(p)
	if err != nil {
		log.Fatal("Conversion failed:", err)
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
