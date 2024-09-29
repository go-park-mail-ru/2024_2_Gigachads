package httpserver

import (
	"log"
	"mail/database"
	"net/http"
)

func checkAuthorizationByID(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "user_id cookie not found", http.StatusForbidden)
		return
	}

	userID := cookie.Value
	log.Printf("User ID from cookie: %s", userID)

	userName, ok := database.UserID[userID]
	if !ok {
		log.Printf("ERROR: user not found in the database for ID: %s", userID)
		http.Error(w, "user not found", http.StatusForbidden)
		return
	}

	log.Printf("User found: %s", userName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authorization successful\n"))
}
