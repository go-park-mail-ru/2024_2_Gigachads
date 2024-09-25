package httpserver

import (
	"github.com/gorilla/mux"
	"log/slog"
	config "mail/config"
	"net/http"
	"slices"
)

type HTTPServer struct {
	server *http.Server
}

func (s *HTTPServer) Start(cfg *config.Config) error {
	s.server = new(http.Server)
	s.server.Addr = cfg.HTTPServer.IP + ":" + cfg.HTTPServer.Port
	s.configureRouter(cfg)
	slog.Info("Server is running on", "port", cfg.HTTPServer.Port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) configureRouter(cfg *config.Config) {
	router := mux.NewRouter()

	// Authorization checking
	router.Use(authMiddleware)

	// CORS checking
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if slices.Contains(cfg.HTTPServer.AllowedIPsByCORS, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
			}
			next.ServeHTTP(w, r)
		})
	})

	router.HandleFunc("/hello", HelloHandler).Methods("GET")
	s.server.Handler = router
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !userIsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// userIsAuthenticated пока замокала, потому что некуда идти проверять валидность - логика аутентификации не написана
func userIsAuthenticated(r *http.Request) bool {
	return true
}
