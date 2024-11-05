package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"GO-API/internal/domain/model"
	"GO-API/internal/pkg/logger"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	logger.Error("Error response: status=%d, message=%s", status, message)

	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    fmt.Sprintf("ERR_%d", status),
	}

	writeJSON(w, status, response)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Failed to encode response: %v", err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	if domainErr, ok := err.(*model.Error); ok {
		switch domainErr.Type {
		case model.ErrorTypeValidation:
			writeError(w, http.StatusBadRequest, domainErr.Message)
		case model.ErrorTypeNotFound:
			writeError(w, http.StatusNotFound, domainErr.Message)
		default:
			writeError(w, http.StatusNotFound, domainErr.Message)
		}
		return
	}
	writeError(w, http.StatusInternalServerError, "internal server error")
}
