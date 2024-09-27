package httpserver

import (
	"net/http"
	"encoding/json"
	//"fmt"
)

type errorResponse struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	response := errorResponse{
		Status: http.StatusForbidden,
		Body:   message,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}