package middleware

import (
	"context"

	"net/http"
)

type contextKey string

const UserIDKey = contextKey("user_id")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), UserIDKey, r.Header.Get("user_id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
