package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDKey = contextKey("User-id")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !userIsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, r.Header.Get(string(UserIDKey)))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func userIsAuthenticated(r *http.Request) bool {
	return true
}
