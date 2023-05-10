package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"shortener/internal/handlers"
)

func StartServer(addr string) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/", handlers.CreateShortHandler)
	r.Get("/{id}", handlers.GetShortHandler)

	fmt.Println("Running server on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
