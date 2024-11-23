package app

import (
	"mail/config"
	"mail/api-service/internal/delivery/httpserver"
	"mail/api-service/pkg/logger"
	"mail/service/postgres"
	"mail/service/redis"
	"mail/api-service/internal/delivery/grpc"
)

func Run(cfg *config.Config, l logger.Logger) error {
	var srv httpserver.HTTPServer

	dbPostgres, err := postgres.Init(cfg)
	if err != nil {
		return err
	}

	redisSessionClient, err := redis.Init(cfg, 0)
	if err != nil {
		return err
	}

	redisCSRFClient, err := redis.Init(cfg, 1)
	if err != nil {
		return err
	}
	clients := grpcClients.Init(cfg)
	
	if err := srv.Start(cfg, dbPostgres, redisSessionClient, redisCSRFClient, clients, l); err != nil {
		return err
	}
	return nil
}
