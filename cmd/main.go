package main

import (
	"flag"
	"log/slog"
	config "mail/config"
	httpserver "mail/internal/delivery/httpserver"
	mw "mail/internal/delivery/middleware"
	repo "mail/internal/repository"
	usecase "mail/internal/usecases"
)

func main() {
	var srv httpserver.HTTPServer
	configPath := flag.String("config-path", "./config/config.yaml", "path to config file")
	flag.Parse()

	config, err := config.GetConfig(*configPath)
	if err != nil {
		slog.Error(err.Error())
	}
	ur := repo.NewUserRepository()
	uu := usecase.NewUserUseCase(ur)
	er := repo.NewEmailRepository()
	eu := usecase.NewEmailUseCase(er)
	sr := repo.NewSessionRepository()
	su := usecase.NewSessionUseCase(sr)
	authHandler := httpserver.NewAuthHandler(uu, su)
	emailHandler := httpserver.NewEmailHandler(eu, su)
	authMW := mw.NewAuthMW(su)
	if err := srv.Start(config, authHandler, emailHandler, authMW); err != nil {
		slog.Error(err.Error())
	}
}
