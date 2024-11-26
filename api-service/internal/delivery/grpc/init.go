package grpcClients

import (
	"fmt"
	"google.golang.org/grpc/credentials/insecure"

	"mail/config"

	"google.golang.org/grpc"
	"mail/gen/go/auth"
	"mail/gen/go/smtp"
)

type Clients struct {
	AuthConn *auth.AuthServiceClient
	SmtpConn *smtp.SmtpPop3ServiceClient
}

func Init(cfg *config.Config) (*Clients, error) {
	//auth microservice
	authConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.AuthServer.IP, cfg.AuthServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf(fmt.Sprintf("the microservice 'authorization' is not available: %v", err))
		return nil, err
	}
	authClient := auth.NewAuthServiceClient(authConn)
	smtpConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.SMTPServer.IP, cfg.SMTPServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf(fmt.Sprintf("the microservice 'authorization' is not available: %v", err))
		return nil, err
	}
	smtpClient := smtp.NewSmtpPop3ServiceClient(smtpConn)
	return &Clients{
		AuthConn: &authClient,
		SmtpConn: &smtpClient,
	}, nil
}
