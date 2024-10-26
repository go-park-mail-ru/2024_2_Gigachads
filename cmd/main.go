package main

import (
	"flag"
	"mail/pkg/logger"
	config "mail/config"
	httpserver "mail/internal/delivery/httpserver"
)

func main() {
	var srv httpserver.HTTPServer

	logger.NewLogger()

	configPath := flag.String("config-path", "./config/config.yaml", "path to config file")
	flag.Parse()

	config, err := config.GetConfig(*configPath)
	if err != nil {
		logger.Error(err.Error())
	}
	if err := srv.Start(config); err != nil {
		logger.Error(err.Error())
	}
}
