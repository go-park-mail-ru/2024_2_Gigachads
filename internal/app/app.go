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

	redisClient, err := redis.Init(cfg)
	if err != nil {
		return err
	}

	if err := srv.Start(cfg, dbPostgres, redisClient, l); err != nil {
		return err
	}
	return nil
}
