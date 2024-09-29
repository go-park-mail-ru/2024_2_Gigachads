package httpserver

import (
	"encoding/json"
	"log/slog"
	"mail/database"
	"net/http"
)

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		response := errorResponse{
			Status: http.StatusForbidden,
			Body:   "Validation_error",
		}
		marshaledResponse, err := json.Marshal(response)
		if err != nil {
			slog.Error("failed to marshal error response")
		}
		w.Write(marshaledResponse)
		return
	}
	userHash := cookie.Value
	
	http.SetCookie(w, &http.Cookie{
		Name:   cookie.Name,
		Value:  "",
		MaxAge: -1,
	})

	delete(database.UserHash, userHash)
	w.WriteHeader(http.StatusOK)
}
