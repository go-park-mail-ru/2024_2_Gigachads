package httpserver

import (
	"github.com/gorilla/mux"
	"log/slog"
	"mail/config"
	"mail/internal/delivery/middleware"
	repo "mail/internal/repository"
	"mail/internal/usecases"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

func (s *HTTPServer) Start(cfg *config.Config) error {
	s.server = new(http.Server)
	s.server.Addr = cfg.HTTPServer.IP + ":" + cfg.HTTPServer.Port
	s.configureRouters(cfg)
	slog.Info("Server is running on", "port", cfg.HTTPServer.Port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) configureRouters(cfg *config.Config) {
	sr := repo.NewSessionRepositoryService()

	ur := repo.NewUserRepositoryService()
	uu := usecase.NewUserService(ur, sr)

	er := repo.NewEmailRepositoryService()
	eu := usecase.NewEmailService(er, sr)

	router := mux.NewRouter()
	public := router.PathPrefix("/").Subrouter()
	private := router.PathPrefix("/").Subrouter()

	ConfigureEmailRouter(private, eu)
	ConfigureAuthRouter(public, private, uu)
	ConfigureAuthMiddleware(private, uu)

	router.Use(func(next http.Handler) http.Handler {
		return middleware.CORS(next, cfg)
	})

	s.server.Handler = router
}
