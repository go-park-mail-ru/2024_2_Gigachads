package usecase

import (
	models "mail/internal/models"
	repository "mail/internal/repository"
)

type SessionUseCase struct {
	repo *repository.SessionRepository
}

func NewSessionUseCase(repo *repository.SessionRepository) *SessionUseCase {
	return &SessionUseCase{
		repo: repo,
	}
}

func (uc *SessionUseCase) CreateSession(mail string) (*models.Session, error) {
	return uc.repo.CreateSession(mail)
}

func (uc *SessionUseCase) GetSession(sessionID string) (*models.Session, error) {
	return uc.repo.GetSession(sessionID)
}

func (uc *SessionUseCase) DeleteSession(sessionID string) error {
	return uc.repo.DeleteSession(sessionID)
}
