package app

import (
	"mail/config"
	"mail/internal/delivery/httpserver"
	"mail/service/postgres"
)

func Run(cfg *config.Config) error {
	var srv httpserver.HTTPServer

	dbPostgres, err := postgres.Init(cfg)
	if err != nil {
		return err
	}

	if err := srv.Start(cfg, dbPostgres); err != nil {
		return err
	}
	return nil
}
