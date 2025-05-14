package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"go_api/internal/service"

	"github.com/go-chi/chi/v5"
)

func GetForecast(w http.ResponseWriter, r *http.Request) {
	log.Println("Servidor iniciado em :8080")
	locationKey := chi.URLParam(r, "locationKey")

	forecast, err := service.GetForecast(locationKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(forecast)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
