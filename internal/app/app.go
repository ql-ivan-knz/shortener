package app

import (
	"log"
	"net/http"
	"shortener/internal/app/routes"
)

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", routes.MainHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
