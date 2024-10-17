package httpserver

import (
	"log/slog"
	config "mail/config"
	"mail/internal/delivery/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	server *http.Server
}

func (s *HTTPServer) Start(cfg *config.Config, authHandler *AuthHandler, emailHandler *EmailHandler, authMW *middleware.AuthMiddleware) error {
	s.server = new(http.Server)
	s.server.Addr = cfg.HTTPServer.IP + ":" + cfg.HTTPServer.Port
	s.configureRouter(cfg, authHandler, emailHandler, authMW)
	slog.Info("Server is running on", "port", cfg.HTTPServer.Port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) configureRouter(cfg *config.Config, authHandler *AuthHandler, emailHandler *EmailHandler, authMW *middleware.AuthMiddleware) {
	router := mux.NewRouter()

	public := router.PathPrefix("/").Subrouter()
	public.HandleFunc("/signup", authHandler.SignUp).Methods("POST", "OPTIONS")
	public.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")

	private := router.PathPrefix("/").Subrouter()
	private.HandleFunc("/mail/inbox", emailHandler.Inbox).Methods("GET", "OPTIONS")
	private.HandleFunc("/logout", authHandler.Logout).Methods("GET", "OPTIONS")
	private.Use(authMW.Handler)

	router.Use(func(next http.Handler) http.Handler {
		return middleware.CORS(next, cfg)
	})

	s.server.Handler = router
}
