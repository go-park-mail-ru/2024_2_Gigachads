package main

import (
	"flag"
	"log/slog"
	config "mail/config"
	app "mail/internal/app"
)

func main() {
	configPath := flag.String("config-path", "./config/config.yaml", "path to config file")
	flag.Parse()

	config, err := config.GetConfig(*configPath)
	if err != nil {
		slog.Error(err.Error())
	}
	if err := app.Run(config); err != nil {
		slog.Error(err.Error())
	}
}
