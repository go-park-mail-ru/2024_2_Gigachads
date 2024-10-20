package utils

import (
	"encoding/json"
	"mail/internal/models"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, r *http.Request, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := models.Error{
		Status: code,
		Body:   message,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}
