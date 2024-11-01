package main

import (
	"flag"
	"mail/pkg/logger"
	config "mail/config"
	httpserver "mail/internal/delivery/httpserver"
)

func main() {
	var srv httpserver.HTTPServer

	l := logger.NewLogger()

	configPath := flag.String("config-path", "./config/config.yaml", "path to config file")
	flag.Parse()

	config, err := config.GetConfig(*configPath, l)
	if err != nil {
		l.Error(err.Error())
	}
	if err := srv.Start(config, l); err != nil {
		l.Error(err.Error())
	}
}
