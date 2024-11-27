package app

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"mail/api-service/pkg/logger"
	repo "mail/auth-service/internal/repo"
	"mail/auth-service/internal/usecase"
	"mail/config"
	proto "mail/gen/go/auth"
	"mail/service/postgres"
	"mail/service/redis"
	"net"
)

func Run(cfg *config.Config, l logger.Logger) error {
	dbPostgres, err := postgres.Init(cfg)
	if err != nil {
		return err
	}
	l.Info("postgres connected")

	ur := repo.NewUserRepositoryService(dbPostgres, l)

	redisSessionClient, err := redis.Init(cfg, 0)
	if err != nil {
		return err
	}
	l.Info("redis session connected")
	sr := repo.NewSessionRepositoryService(redisSessionClient, l)

	redisCSRFClient, err := redis.Init(cfg, 1)
	if err != nil {
		return err
	}
	l.Info("redis csrf connected")
	cr := repo.NewCsrfRepositoryService(redisCSRFClient, l)

	port := ":" + cfg.AuthServer.Port
	conn, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	l.Info("auth microservice started")

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(),
	))

	proto.RegisterAuthServiceServer(server, usecase.NewAuthServer(ur, sr, cr))

	err = server.Serve(conn)
	if err != nil {
		return err
	}
	return nil
}
