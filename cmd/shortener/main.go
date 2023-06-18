package main

import (
	"context"
	"log"
	"shortener/config"
	"shortener/internal/server"
	"shortener/internal/storage"
)

func main() {
	cfg := config.GetConfig()
	s, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	h := server.NewHandlers(cfg, s, context.Background())

	server.StartServer(h)
}
