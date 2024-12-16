package httpserver

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
	"github.com/gorilla/mux"
	"mail/api-service/internal/delivery/grpc"
	authRouter "mail/api-service/internal/delivery/httpserver/auth"
	emailRouter "mail/api-service/internal/delivery/httpserver/email"
	userRouter "mail/api-service/internal/delivery/httpserver/user"
	mw "mail/api-service/internal/delivery/middleware"
	repo "mail/api-service/internal/repository"
	usecase "mail/api-service/internal/usecases"
	"mail/api-service/pkg/logger"
	"mail/config"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

func (s *HTTPServer) Start(cfg *config.Config, db *sql.DB, clients *grpcClients.Clients, redisLastModifiedClient *redis.Client, l logger.Logable) error {
	s.server = new(http.Server)
	s.server.Addr = cfg.HTTPServer.IP + ":" + cfg.HTTPServer.Port
	s.configureRouters(cfg, db, clients, redisLastModifiedClient, l)
	l.Info("Server is running on", "port", cfg.HTTPServer.Port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) configureRouters(cfg *config.Config, db *sql.DB, clients *grpcClients.Clients, redisLastModifiedClient *redis.Client, l logger.Logable) {

	ur := repo.NewUserRepositoryService(db)
	uu := usecase.NewUserService(ur)
	au := usecase.NewAuthService(*clients.AuthConn)
	er := repo.NewEmailRepositoryService(db, redisLastModifiedClient, l)
	eu := usecase.NewEmailService(er, *clients.SmtpConn)

	router := mux.NewRouter()
	router = router.PathPrefix("/api/").Subrouter()
	router.Use(mw.PanicMiddleware)
	router.Use(mw.NewLogMW(l).Handler)

	authRout := authRouter.NewAuthRouter(au)
	emailRout := emailRouter.NewEmailRouter(eu)
	userRout := userRouter.NewUserRouter(uu)
	mwAuth := mw.NewAuthMW(au)

	emailRout.ConfigureEmailRouter(router)
	authRout.ConfigureAuthRouter(router)
	userRout.ConfigureUserRouter(router)

	handler := mw.ConfigureMWs(cfg, router, mwAuth)

	s.server.Handler = handler

}
