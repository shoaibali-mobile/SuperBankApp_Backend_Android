package handlers

import (
	"encoding/json"
	"net/http"

	"bankapp-microservices/internal/models"
)

func respondWithSuccess(w http.ResponseWriter, data interface{}, message ...string) {
	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		Success: true,
		Data:    data,
	}
	if len(message) > 0 {
		response.Message = message[0]
	}
	json.NewEncoder(w).Encode(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.Response{
		Success: false,
		Message: message,
	})
}
