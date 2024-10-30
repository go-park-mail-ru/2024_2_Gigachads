package httpserver

import (
	"database/sql"
	"log/slog"
	"mail/config"
	authRouter "mail/internal/delivery/httpserver/auth"
	emailRouter "mail/internal/delivery/httpserver/email"
	mw "mail/internal/delivery/middleware"
	repo "mail/internal/repository"
	usecase "mail/internal/usecases"
	"mail/pkg/smtp"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	server *http.Server
}

func (s *HTTPServer) Start(cfg *config.Config, db *sql.DB) error {
	s.server = new(http.Server)
	s.server.Addr = cfg.HTTPServer.IP + ":" + cfg.HTTPServer.Port
	s.configureRouters(cfg, db)
	slog.Info("Server is running on", "port", cfg.HTTPServer.Port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) configureRouters(cfg *config.Config, db *sql.DB) {
	sr := repo.NewSessionRepositoryService()

	smtpClient := s.createAndConfigureSMTPClient(cfg)

	ur := repo.NewUserRepositoryService(db)
	smtpRepo := repo.NewSMTPRepository(smtpClient, cfg)
	uu := usecase.NewUserService(ur, sr)

	er := repo.NewEmailRepositoryService(db)
	eu := usecase.NewEmailService(er, sr, smtpRepo)

	router := mux.NewRouter()
	router = router.PathPrefix("/").Subrouter()

	authRout := authRouter.NewAuthRouter(uu)
	emailRout := emailRouter.NewEmailRouter(eu)
	mwAuth := mw.NewAuthMW(uu)

	emailRout.ConfigureEmailRouter(router)
	authRout.ConfigureAuthRouter(router)

	handler := mw.ConfigureMWs(cfg, router, mwAuth)

	s.server.Handler = handler
}

func (s *HTTPServer) createAndConfigureSMTPClient(cfg *config.Config) *smtp.SMTPClient {
	return smtp.NewSMTPClient(
		cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password,
	)
}
