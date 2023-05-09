package server

import (
	"log"
	"net/http"
	"shortener/internal/handlers"
)

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.MainHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
