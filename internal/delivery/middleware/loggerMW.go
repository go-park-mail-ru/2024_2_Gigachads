package middleware

import (
	"context"
	"net/http" 
	"mail/pkg/logger"
	"mail/pkg/utils"
)

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *LogResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.RequestURI
		id, ok := r.Context().Value("requestID").(string)
		if !ok {
			id, _ = utils.GenerateHash()
		}
		
		logger.Info("User entered", "url", path, "method", method, "requestID", id)

		logRW := &LogResponseWriter{w, http.StatusOK}
		ctx := context.WithValue(r.Context(), "requestID", id)

		next.ServeHTTP(logRW, r.WithContext(ctx))

		statusCode := logRW.statusCode
		logger.Info("User left", "url", path, "status code", statusCode, "requestID", id)
	})
}