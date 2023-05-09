package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"shortener/internal/handlers"
)

func StartServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/", handlers.CreateShortHandler)
	r.Get("/{id}", handlers.GetShortHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
