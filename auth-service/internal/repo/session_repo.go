package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"mail/api-service/pkg/logger"
	"mail/api-service/pkg/utils"
	"mail/models"
	"time"
)

type SessionRepositoryService struct {
	repo   *redis.Client
	logger logger.Logable
}

func NewSessionRepositoryService(client *redis.Client, l logger.Logable) models.SessionRepository {
	return &SessionRepositoryService{repo: client, logger: l}
}

func (sr *SessionRepositoryService) CreateSession(ctx context.Context, mail string) (*models.Session, error) {

	mail = utils.Sanitize(mail)

	hash, err := utils.GenerateHash()
	if err != nil {
		sr.logger.Error(err.Error())
		return &models.Session{}, err
	}
	expiration := time.Now().Add(24 * time.Hour)

	err = sr.repo.Set(ctx, string(hash), mail, 24*time.Hour).Err()
	if err != nil {
		sr.logger.Error(err.Error())
		return nil, err
	}

	hash = utils.Sanitize(hash)

	session := models.Session{Name: "email", ID: hash, Time: expiration, UserLogin: mail}

	return &session, nil
}

func (sr *SessionRepositoryService) DeleteSession(ctx context.Context, sessionID string) error {

	sessionID = utils.Sanitize(sessionID)

	err := sr.repo.Del(ctx, sessionID).Err()
	if err != nil {
		sr.logger.Error(err.Error())
		return err
	}
	return nil
}

func (sr *SessionRepositoryService) GetSession(ctx context.Context, sessionID string) (string, error) {

	sessionID = utils.Sanitize(sessionID)

	email, err := sr.repo.Get(ctx, sessionID).Result()
	if err != nil {
		sr.logger.Error(err.Error())
		return "", err
	}

	email = utils.Sanitize(email)

	return email, nil
}
