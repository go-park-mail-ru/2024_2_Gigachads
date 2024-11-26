package main

import (
	"flag"
	"mail/api-service/pkg/logger"
	"mail/config"
	"mail/smtp-service/internal/app"
)

func main() {
	l := logger.NewLogger()
	configPath := flag.String("config-path", "./config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.GetConfig(*configPath, l)
	if err != nil {
		l.Error(err.Error())
	}
	if err := app.Run(cfg, l); err != nil {
		l.Error(err.Error())
	}

}
