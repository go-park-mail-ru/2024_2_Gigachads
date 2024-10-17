package usecase

import (
	"mail/internal/repository"
)

type SMTPUsecase struct {
	smtpRepository *repository.SMTPRepository
}

func NewSMTPUsecase(smtpRepository *repository.SMTPRepository) *SMTPUsecase {
	return &SMTPUsecase{smtpRepository: smtpRepository}
}

func (s *SMTPUsecase) SendEmail(from string, to []string, subject string, body string) error {
	return s.smtpRepository.SendEmail(from, to, subject, body)
}
