package repository

import (
	"context"
	models "mail/internal/models"
	"mail/pkg/utils"
	"time"
	"github.com/redis/go-redis/v9"
)

type SessionRepositoryService struct {
	repo *redis.Client
}

func NewSessionRepositoryService(client *redis.Client) models.SessionRepository {
	return &SessionRepositoryService{repo: client}
}

func (sr *SessionRepositoryService) CreateSession(ctx context.Context, mail string) (*models.Session, error) {
	hash, err := utils.GenerateHash()
	if err != nil {
		return &models.Session{}, err
	}
	expiration := time.Now().Add(24 * time.Hour)
	
	err = sr.repo.Set(ctx, string(hash), mail, 24 * time.Hour).Err()
	if err != nil {
		return nil, err
	}
	session := models.Session{Name: "email", ID: hash, Time: expiration, UserLogin: mail}
	
	return &session, nil
}

func (sr *SessionRepositoryService) DeleteSession(ctx context.Context, sessionID string) error {
	err := sr.repo.Del(ctx, sessionID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (sr *SessionRepositoryService) GetSession(ctx context.Context, sessionID string) (string, error) {
	email, err := sr.repo.Get(ctx, sessionID).Result()
	if err != nil {
		return "", err
	}
	return email, nil
}
