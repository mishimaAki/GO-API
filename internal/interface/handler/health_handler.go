package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"GO-API/internal/pkg/logger"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"time"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	logger.Info("Health check requested")

	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logger.Debug("Health check responded successfully")
}
