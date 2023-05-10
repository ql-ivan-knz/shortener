package main

import (
	"shortener/config"
	"shortener/internal/server"
)

func main() {
	config.ParseFlags()

	server.StartServer(config.Addr)
}
