package app

import (
	"mail/api-service/internal/delivery/grpc"
	"mail/api-service/internal/delivery/httpserver"
	"mail/api-service/pkg/logger"
	"mail/config"
	"mail/service/postgres"
)

func Run(cfg *config.Config, l logger.Logger) error {
	var srv httpserver.HTTPServer

	dbPostgres, err := postgres.Init(cfg)
	if err != nil {
		return err
	}
	l.Info("postgres connected")

	clients, err := grpcClients.Init(cfg, l)
	if err != nil {
		return err
	}

	if err := srv.Start(cfg, dbPostgres, clients, l); err != nil {
		return err
	}
	return nil
}
