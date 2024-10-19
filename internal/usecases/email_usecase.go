package usecase

import (
	models "mail/internal/models"
)

type EmailService struct {
	EmailRepo   models.EmailRepository
	SessionRepo models.SessionRepository
}

func NewEmailService(erepo models.EmailRepository, srepo models.SessionRepository) models.EmailUseCase {
	return &EmailService{
		EmailRepo:   erepo,
		SessionRepo: srepo,
	}
}

func (es *EmailService) Inbox(sessionID string) ([]models.Email, error) {
	session, err := es.SessionRepo.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	return es.EmailRepo.Inbox(session.UserLogin)
}
