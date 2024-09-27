package middleware

import (
	"mail/config"
	"net/http"
)

// func CORS(next http.Handler, cfg *config.Config) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		origin := r.Header.Get("Origin")
// 		if slices.Contains(cfg.HTTPServer.AllowedIPsByCORS, origin) {
// 			w.Header().Set("Access-Control-Allow-Origin", origin)
// 			w.Header().Add("Vary", "Origin")
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

func CORS(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", cfg.HTTPServer.AllowedIPsByCORS[0])
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", cfg.HTTPServer.AllowedIPsByCORS[0])
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(w, r)
	})
}
