package main

import (
	"flag"
	"log/slog"
	config "mail/config"
	httpserver "mail/internal/app/httpserver"
)

func main() {
	var srv httpserver.HTTPServer
	configPath := flag.String("config-path", "./config/config.yaml", "path to config file")

	config, err := config.GetConfig(*configPath)
	if err != nil {
		slog.Error(err.Error())
	}

	if err := srv.Start(config); err != nil {
		slog.Error(err.Error())
	}
}
