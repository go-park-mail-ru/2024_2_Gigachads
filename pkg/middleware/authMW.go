package middleware

import (
	"context"
	"mail/database"
	"net/http"
)

type contextKey string

const UserIDKey = contextKey("user_id")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Пропуск аутентификации для предзапросов CORS
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("user_id")
		if err != nil || !checkAuthorization(cookie.Value) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkAuthorization(userID string) bool {
	_, exists := database.UserID[userID]
	return exists
}
