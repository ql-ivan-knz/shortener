package config

import (
	"flag"
	"os"
)

var (
	ServerAddr string
	BaseURL    string
)

func ParseFlags() {
	flag.StringVar(&ServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "")

	flag.Parse()

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		ServerAddr = envAddr
	}

	if resAddr := os.Getenv("BASE_URL"); resAddr != "" {
		BaseURL = resAddr
	}
}
