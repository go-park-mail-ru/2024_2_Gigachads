package usecase

import (
	"context"
	"fmt"
	proto "mail/gen/go/smtp"
	"mail/smtp-service/internal/models"
	"mail/smtp-service/pkg/pop3"
	"mail/smtp-service/pkg/smtp"
	"time"
)

type SmtpPop3ServiceServer struct {
	proto.UnimplementedSmtpPop3ServiceServer
	POP3Repo   models.POP3Repository
	SMTPRepo   models.SMTPRepository
	EmailRepo  models.EmailRepositorySMTP
	pop3Client *pop3.Pop3Client
	smtpClient *smtp.SMTPClient
}

func NewSmtpPop3ServiceServer(prepo models.POP3Repository, srepo models.SMTPRepository, erepo models.EmailRepositorySMTP, pop3Client *pop3.Pop3Client, smtpClient *smtp.SMTPClient) proto.SmtpPop3ServiceServer {
	return &SmtpPop3ServiceServer{
		POP3Repo:   prepo,
		SMTPRepo:   srepo,
		EmailRepo:  erepo,
		pop3Client: pop3Client,
		smtpClient: smtpClient,
	}
}

func (s *SmtpPop3ServiceServer) SendEmail(ctx context.Context, request *proto.SendEmailRequest) (*proto.SendEmailReply, error) {
	to := []string{request.GetTo()}
	return &proto.SendEmailReply{}, s.SMTPRepo.SendEmail(request.GetFrom(), to, request.GetSubject(), request.GetBody())
}

func (s *SmtpPop3ServiceServer) ForwardEmail(ctx context.Context, request *proto.ForwardEmailRequest) (*proto.ForwardEmailReply, error) {
	forwardSubject := "Fwd: " + request.GetTitle()
	forwardBody := fmt.Sprintf(`
---------- Forwarded message ---------
From: %s
Date: %s
Subject: %s

%s
`, request.Sender, request.GetSendingDate().AsTime().Format(time.RFC1123),
		request.GetTitle(), request.GetDescription())
	to := []string{request.GetTo()}
	return &proto.ForwardEmailReply{}, s.SMTPRepo.SendEmail(request.GetFrom(), to, forwardSubject, forwardBody)
}

func (s *SmtpPop3ServiceServer) ReplyEmail(ctx context.Context, request *proto.ReplyEmailRequest) (*proto.ReplyEmailReply, error) {
	replySubject := "Re: " + request.GetTitle()
	replyBody := fmt.Sprintf(`
%s

On %s, %s wrote:
> %s
`, request.GetReplyText(), request.GetSendingDate().AsTime().Format(time.RFC1123),
		request.GetSender(), request.GetDescription())
	to := []string{request.GetTo()}
	return &proto.ReplyEmailReply{}, s.SMTPRepo.SendEmail(request.GetFrom(), to, replySubject, replyBody)
}

func (s *SmtpPop3ServiceServer) FetchEmailsViaPOP3(ctx context.Context, request *proto.FetchEmailsViaPOP3Request) (*proto.FetchEmailsViaPOP3Reply, error) {
	if err := s.POP3Repo.Connect(); err != nil {
		return &proto.FetchEmailsViaPOP3Reply{}, err
	}
	defer s.POP3Repo.Quit()

	if err := s.POP3Repo.FetchEmails(s.EmailRepo); err != nil {
		return &proto.FetchEmailsViaPOP3Reply{}, err
	}
	return &proto.FetchEmailsViaPOP3Reply{}, nil
}
