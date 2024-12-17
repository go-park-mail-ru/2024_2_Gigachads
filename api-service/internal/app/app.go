package app

import (
	"mail/api-service/internal/delivery/grpc"
	"mail/api-service/internal/delivery/httpserver"
	"mail/api-service/pkg/logger"
	"mail/config"
	"mail/service/postgres"
	"mail/service/redis"
)

func Run(cfg *config.Config, l logger.Logger) error {
	var srv httpserver.HTTPServer

	dbPostgres, err := postgres.Init(cfg)
	if err != nil {
		return err
	}
	l.Info("postgres connected")

	redisLastModifiedClient, err := redis.Init(cfg, 2)
	if err != nil {
		return err
	}
	l.Info("lastmodified redis connected")

	clients, err := grpcClients.Init(cfg, l)
	if err != nil {
		return err
	}

	if err := srv.Start(cfg, dbPostgres, clients, redisLastModifiedClient, l); err != nil {
		return err
	}
	return nil
}
