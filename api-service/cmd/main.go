package main

import (
	"flag"
	config "mail/config"
	app "mail/api-service/internal/app"
	"mail/api-service/pkg/logger"
)

func main() {
	l := logger.NewLogger()

	configPath := flag.String("config-path", "./config/config.yaml", "path to config file")
	flag.Parse()

	config, err := config.GetConfig(*configPath, l)
	if err != nil {
		l.Error(err.Error())
	}
	if err := app.Run(config, l); err != nil {
		l.Error(err.Error())
	}
}
