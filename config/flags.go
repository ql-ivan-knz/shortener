package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddr string
	BaseURL    string
}

func GetConfig() Config {
	var cfg = Config{}

	flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "")

	flag.Parse()

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		cfg.ServerAddr = envAddr
	}

	if resAddr := os.Getenv("BASE_URL"); resAddr != "" {
		cfg.BaseURL = resAddr
	}

	return cfg
}
