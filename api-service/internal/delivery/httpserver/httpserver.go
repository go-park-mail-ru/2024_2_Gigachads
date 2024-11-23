package httpserver

import (
	"mail/config"
	authRouter "mail/api-service/internal/delivery/httpserver/auth"
	emailRouter "mail/api-service/internal/delivery/httpserver/email"
	userRouter "mail/api-service/internal/delivery/httpserver/user"
	mw "mail/api-service/internal/delivery/middleware"
	"mail/api-service/internal/models"
	repo "mail/api-service/internal/repository"
	usecase "mail/api-service/internal/usecases"
	"mail/api-service/pkg/logger"
	"mail/api-service/pkg/pop3"
	"mail/api-service/pkg/smtp"
	
	"database/sql"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type HTTPServer struct {
	server *http.Server
}

func (s *HTTPServer) Start(cfg *config.Config, db *sql.DB, redisSession *redis.Client, redisCSRF *redis.Client, l logger.Logable) error {
	s.server = new(http.Server)
	s.server.Addr = cfg.HTTPServer.IP + ":" + cfg.HTTPServer.Port
	s.configureRouters(cfg, db, redisSession, redisCSRF, l)
	l.Info("Server is running on", "port", cfg.HTTPServer.Port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) configureRouters(cfg *config.Config, db *sql.DB, redisSession *redis.Client, redisCSRF *redis.Client, clients *grpcClients.Clients, l logger.Logable) {
	
	sr := repo.NewSessionRepositoryService(redisSession, l)
	cr := repo.NewCsrfRepositoryService(redisCSRF, l)
	smtpClient := s.createAndConfigureSMTPClient(cfg)

	ur := repo.NewUserRepositoryService(db)
	smtpRepo := repo.NewSMTPRepository(smtpClient, cfg)
	uu := usecase.NewUserService(ur, sr, cr)
	au := usecase.NewAuthService(clients.AuthConn)
	er := repo.NewEmailRepositoryService(db, l)

	pop3Client := s.createAndConfigurePOP3Client(cfg)

	eu := usecase.NewEmailService(er, sr, smtpRepo, pop3Client)

	router := mux.NewRouter()
	router = router.PathPrefix("/").Subrouter()
	router.Use(mw.PanicMiddleware)
	router.Use(mw.NewLogMW(l).Handler)

	authRout := authRouter.NewAuthRouter(uu)
	emailRout := emailRouter.NewEmailRouter(eu)
	userRout := userRouter.NewUserRouter(uu)
	mwAuth := mw.NewAuthMW(uu)

	emailRout.ConfigureEmailRouter(router)
	authRout.ConfigureAuthRouter(router)
	userRout.ConfigureUserRouter(router)

	handler := mw.ConfigureMWs(cfg, router, mwAuth)

	s.server.Handler = handler

	s.startEmailFetcher(eu)
}

func (s *HTTPServer) createAndConfigureSMTPClient(cfg *config.Config) *smtp.SMTPClient {
	return smtp.NewSMTPClient(
		cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password,
	)
}
func (s *HTTPServer) createAndConfigurePOP3Client(cfg *config.Config) *pop3.Pop3Client {
	return pop3.NewPop3Client(cfg.Pop3.Host, cfg.Pop3.Port, cfg.Pop3.Username, cfg.Pop3.Password)
}

func (s *HTTPServer) startEmailFetcher(eu models.EmailUseCase) {
	fetcher := email.NewEmailFetcher(eu)
	fetcher.Start()
}
