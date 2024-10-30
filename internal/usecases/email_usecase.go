package usecase

import (
	models "mail/internal/models"
)

type EmailService struct {
	EmailRepo   models.EmailRepository
	SessionRepo models.SessionRepository
	SMTPRepo    models.SMTPRepository
}

func NewEmailService(erepo models.EmailRepository, srepo models.SessionRepository, smtprepo models.SMTPRepository) *EmailService {
	return &EmailService{
		EmailRepo:   erepo,
		SessionRepo: srepo,
		SMTPRepo:    smtprepo,
	}
}

func (es *EmailService) Inbox(sessionID string) ([]models.Email, error) {
	session, err := es.SessionRepo.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	return es.EmailRepo.Inbox(session.UserLogin)
}

func (es *EmailService) SendEmail(from string, to []string, subject string, body string) error {
	return es.SMTPRepo.SendEmail(from, to, subject, body)
}
