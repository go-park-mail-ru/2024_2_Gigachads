package app

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"mail/api-service/pkg/logger"
	"mail/config"
	proto "mail/gen/go/smtp"
	"mail/service/postgres"
	"mail/smtp-service/internal/delivery"
	"mail/smtp-service/internal/repo"
	"mail/smtp-service/internal/usecase"
	"mail/smtp-service/pkg/pop3"
	"mail/smtp-service/pkg/smtp"
	"net"
)

func Run(cfg *config.Config, l logger.Logger) error {
	dbPostgres, err := postgres.Init(cfg)
	if err != nil {
		return err
	}

	er := repo.NewEmailRepositoryService(dbPostgres, l)

	smtpClient := createAndConfigureSMTPClient(cfg)
	smtpRepo := repo.NewSMTPRepository(smtpClient, cfg)

	pop3Client := createAndConfigurePOP3Client(cfg)
	pop3Repo := repo.NewPOP3Repository(pop3Client, cfg)

	delivery.StartEmailFetcher(pop3Repo, er)

	port := ":" + cfg.SMTPServer.Port
	conn, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	fmt.Println("smtp started")

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(),
	))

	proto.RegisterSmtpPop3ServiceServer(server, usecase.NewSmtpPop3ServiceServer(pop3Repo, smtpRepo, er, pop3Client, smtpClient))

	err = server.Serve(conn)
	if err != nil {
		return err
	}
	return nil
}

func createAndConfigureSMTPClient(cfg *config.Config) *smtp.SMTPClient {
	return smtp.NewSMTPClient(
		cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password,
	)
}
func createAndConfigurePOP3Client(cfg *config.Config) *pop3.Pop3Client {
	return pop3.NewPop3Client(cfg.Pop3.Host, cfg.Pop3.Port, cfg.Pop3.Username, cfg.Pop3.Password)
}
