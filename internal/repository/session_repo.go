package repository

import (
	"fmt"
	models "mail/internal/models"
	"mail/pkg/utils"
	"time"
)

type SessionRepository interface {
	CreateSession(mail string) (*models.Session, error)
	DeleteSession(sessionID string) error
	GetSession(sessionID string) (*models.Session, error)
}

type SessionRepositoryService struct {
	repo map[string]*models.Session
}

func NewSessionRepositoryService() SessionRepository {
	repo := make(map[string]*models.Session)
	return &SessionRepositoryService{repo: repo}
}

func (sr *SessionRepositoryService) CreateSession(mail string) (*models.Session, error) {
	hash, err := utils.GenerateHash()
	if err != nil {
		return &models.Session{}, err
	}
	expiration := time.Now().Add(24 * time.Hour)
	session := models.Session{Name: "session", ID: hash, Time: expiration, UserLogin: mail}
	sr.repo[mail] = &session
	return &session, nil
}

func (sr *SessionRepositoryService) DeleteSession(sessionID string) error {
	delete(sr.repo, sessionID)
	return nil
}

func (sr *SessionRepositoryService) GetSession(sessionID string) (*models.Session, error) {
	session, ok := sr.repo[sessionID]
	if !ok {
		return &models.Session{}, fmt.Errorf("not_auth")
	}
	return session, nil
}
