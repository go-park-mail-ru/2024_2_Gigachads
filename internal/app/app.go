package app

import (
	"mail/config"
	"mail/internal/delivery/httpserver"
	"mail/pkg/logger"
	"mail/service/postgres"
	"mail/service/redis"
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

	if err := srv.Start(cfg, dbPostgres, redisSessionClient, redisCSRFClient, l); err != nil {
		return err
	}
	return nil
}
