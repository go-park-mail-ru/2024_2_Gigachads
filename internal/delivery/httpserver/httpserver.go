package httpserver

import (
	"github.com/gorilla/mux"
	"mail/pkg/logger"
	"mail/config"
	authRouter "mail/internal/delivery/httpserver/auth"
	emailRouter "mail/internal/delivery/httpserver/email"
	mw "mail/internal/delivery/middleware"
	repo "mail/internal/repository"
	"mail/internal/usecases"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

func (s *HTTPServer) Start(cfg *config.Config, l logger.Logable) error {
	s.server = new(http.Server)
	s.server.Addr = cfg.HTTPServer.IP + ":" + cfg.HTTPServer.Port
	//l := logger.NewLogger()
	s.configureRouters(cfg, l)
	l.Info("Server is running on", "port", cfg.HTTPServer.Port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) configureRouters(cfg *config.Config, l logger.Logable) {
	sr := repo.NewSessionRepositoryService()

	ur := repo.NewUserRepositoryService()
	uu := usecase.NewUserService(ur, sr)

	er := repo.NewEmailRepositoryService()
	eu := usecase.NewEmailService(er, sr)

	router := mux.NewRouter()
	router = router.PathPrefix("/").Subrouter()
	router.Use(mw.PanicMiddleware)
	router.Use(mw.NewLogMW(l).Handler)

	authRout := authRouter.NewAuthRouter(uu)
	emailRout := emailRouter.NewEmailRouter(eu)
	mwAuth := mw.NewAuthMW(uu)

	emailRout.ConfigureEmailRouter(router)
	authRout.ConfigureAuthRouter(router)

	handler := mw.ConfigureMWs(cfg, router, mwAuth)

	s.server.Handler = handler
}
