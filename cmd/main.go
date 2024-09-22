package main

import (
	httpserver "mail/internal/app/httpserver"
	config "mail/config"
)

func main() {
	var srv httpserver.HTTPServer
	config := config.GetConfig("../config/config.yaml")
	srv.Start(config)
}