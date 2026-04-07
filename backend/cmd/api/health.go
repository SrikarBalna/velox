package main

import (
	"encoding/json"
	"net/http"
)

// healthHandler reports the API status.
// @Summary Health Check
// @Description Check if the API server is healthy.
// @Tags General
// @Produce json
// @Success 200 {object} map[string]string "Healthy"
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
