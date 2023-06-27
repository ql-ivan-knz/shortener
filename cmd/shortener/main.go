package main

import (
	"log"
	"shortener/config"
	"shortener/internal/middleware/logger"
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

	lg, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
		return
	}

	m := server.NewMiddleware(lg)
	h := server.NewHandlers(cfg, s, lg)

	err = server.Run(h, m)
	if err != nil {
		log.Fatal(err)
		return
	}
}
