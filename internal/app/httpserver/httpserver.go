package httpserver

import (
	"github.com/gorilla/mux"
	"log/slog"
	config "mail/config"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

func (s *HTTPServer) Start(config *config.Config) error {
	s.server = new(http.Server)
	s.server.Addr = config.HTTPServer.IP + ":" + config.HTTPServer.Port
	s.configureRouter()
	slog.Info("Server is running on", "port", config.HTTPServer.Port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) configureRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/hello", HelloHandler).Methods("GET")
	router.HandleFunc("/signup", SignUpHandler).Methods("POST", "GET", "OPTIONS")
	router.HandleFunc("/login", LogInHandler).Methods("POST", "GET", "OPTIONS")
	s.server.Handler = router
}
