package usecase

import (
	"fmt"
	proto "mail/gen/go/smtp"
	"mail/models"
)

type SmtpPop3Server struct {
	proto.UnimplementedSmtpPop3Server
	UserRepo    models.UserRepository
	SessionRepo models.SessionRepository
	CsrfRepo    models.CsrfRepository
}

func NewSmtpPop3Server(urepo models.UserRepository, srepo models.SessionRepository, crepo models.CsrfRepository) proto.SmtpPop3Server {
	return &SmtpPop3Server{
		UserRepo:    urepo,
		SessionRepo: srepo,
		CsrfRepo:    crepo,
	}
}

func (as *SmtpPop3Server) Signup(ctx context.Context, signup *proto.SignupRequest) (*proto.SignupReply, error) {
	
}
