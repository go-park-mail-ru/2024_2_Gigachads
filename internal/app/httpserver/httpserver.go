package httpserver

import (
	"github.com/gorilla/mux"
	"log/slog"
	config "mail/config"
	"mail/pkg/middleware"
	"net/http"
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

	router.Use(middleware.AuthMiddleware)
	router.Use(func(next http.Handler) http.Handler {
		return middleware.CORS(next, cfg)
	})

	router.HandleFunc("/hello", HelloHandler).Methods("GET")
	s.server.Handler = router
}
