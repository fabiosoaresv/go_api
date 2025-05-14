package main

import (
	"log"
	"net/http"

	"go_api/internal/handler"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/forecast/{locationKey}", handler.GetForecast)

	log.Println("Servidor iniciado em :8080")
	http.ListenAndServe(":8080", r)
}
