package repository

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mail/internal/models"
	"time"
)

type SessionRepository struct {
	repo map[string]*models.Session
}

func NewSessionRepository() *SessionRepository {
	repo := make(map[string]*models.Session)
	return &SessionRepository{repo: repo}
}

func (sr *SessionRepository) CreateSession(mail string) (*models.Session, error) {
	hash, err := GenerateHash()
	if err != nil {
		return &models.Session{}, err
	}
	expiration := time.Now().Add(24 * time.Hour)
	session := models.Session{Name: "session", ID: hash, Time: expiration, UserLogin: mail}
	sr.repo[mail] = &session
	return &session, nil
}

func (sr *SessionRepository) DeleteSession(sessionID string) error {
	delete(sr.repo, sessionID)
	return nil
}

func (sr *SessionRepository) GetSession(sessionID string) (*models.Session, error) {
	session, ok := sr.repo[sessionID]
	if !ok {
		return &models.Session{}, fmt.Errorf("not_auth")
	}
	return session, nil
}

func GenerateHash() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return hex.EncodeToString(bytes), nil
}
