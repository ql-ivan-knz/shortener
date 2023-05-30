package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddr      string
	BaseURL         string
	FileStoragePath string
}

func GetConfig() Config {
	var cfg = Config{}

	flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "")
	flag.StringVar(&cfg.FileStoragePath, "f", "short-url-db.json", "file storage path")

	flag.Parse()

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		cfg.ServerAddr = envAddr
	}

	if resAddr := os.Getenv("BASE_URL"); resAddr != "" {
		cfg.BaseURL = resAddr
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}

	return cfg
}
