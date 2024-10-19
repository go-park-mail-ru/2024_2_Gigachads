package usecase

import (
	models "mail/internal/models"
	repository "mail/internal/repository"
)

type EmailUseCase interface {
	Inbox(id string) ([]models.Email, error)
}

type EmailService struct {
	EmailRepo   repository.EmailRepository
	SessionRepo repository.SessionRepository
}

func NewEmailService(erepo repository.EmailRepository, srepo repository.SessionRepository) EmailUseCase {
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
