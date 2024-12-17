package middleware

import (
	"mail/config"
	"net/http"
)

func CORS(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Проверяем, является ли origin разрешенным
		allowedOrigin := ""
		for _, allowed := range cfg.HTTPServer.AllowedIPsByCORS {
			if allowed == origin {
				allowedOrigin = origin
				break
			}
		}

		// Устанавливаем заголовок только если origin разрешен
		if allowedOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
