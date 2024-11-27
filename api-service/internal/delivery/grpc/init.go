package grpcClients

import (
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"mail/api-service/pkg/logger"

	"mail/config"

	"google.golang.org/grpc"
	"mail/gen/go/auth"
	"mail/gen/go/smtp"
)

type Clients struct {
	AuthConn *auth.AuthServiceClient
	SmtpConn *smtp.SmtpPop3ServiceClient
}

func Init(cfg *config.Config, l logger.Logable) (*Clients, error) {
	//auth microservice
	authConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.AuthServer.IP, cfg.AuthServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Error("auth is not available")
		return nil, err
	}
	authClient := auth.NewAuthServiceClient(authConn)
	l.Info("auth connected successfully")

	smtpConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.SMTPServer.IP, cfg.SMTPServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Error("smtp is not available")
		return nil, err
	}
	smtpClient := smtp.NewSmtpPop3ServiceClient(smtpConn)
	l.Info("smtp connected successfully")

	return &Clients{
		AuthConn: &authClient,
		SmtpConn: &smtpClient,
	}, nil
}
