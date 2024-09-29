package middleware

import (
	"context"
	"mail/database"
	"net/http"
)

type contextKey string

const Key = contextKey("session")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Пропуск аутентификации для предзапросов CORS
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("session")
		if err != nil || !checkAuthorization(cookie.Value) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), Key, cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkAuthorization(hash string) bool {
	_, exists := database.UserHash[hash]
	return exists
}
