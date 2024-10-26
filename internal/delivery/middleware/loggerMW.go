package middleware

import (
	"net/http"
	"mail/pkg/logger"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.RequestURI
		logger.Info("User entered", "url", path, "method", method)
		next.ServeHTTP(w, r)
	})
}