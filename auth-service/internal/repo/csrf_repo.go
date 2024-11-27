package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"mail/api-service/pkg/logger"
	"mail/api-service/pkg/utils"
	"mail/auth-service/internal/models"
	"time"
)

type CsrfRepositoryService struct {
	repo   *redis.Client
	logger logger.Logable
}

func NewCsrfRepositoryService(client *redis.Client, l logger.Logable) models.CsrfRepository {
	return &CsrfRepositoryService{repo: client, logger: l}
}

func (sr *CsrfRepositoryService) CreateCsrf(ctx context.Context, mail string) (*models.Csrf, error) {

	mail = utils.Sanitize(mail)

	hash, err := utils.GenerateHash()
	if err != nil {
		sr.logger.Error(err.Error())
		return &models.Csrf{}, err
	}
	expiration := time.Now().Add(24 * time.Hour)

	err = sr.repo.Set(ctx, string(hash), mail, 24*time.Hour).Err()
	if err != nil {
		sr.logger.Error(err.Error())
		return nil, err
	}

	hash = utils.Sanitize(hash)

	Csrf := models.Csrf{Name: "csrf", ID: hash, Time: expiration, UserLogin: mail}

	return &Csrf, nil
}

func (sr *CsrfRepositoryService) DeleteCsrf(ctx context.Context, CsrfID string) error {

	CsrfID = utils.Sanitize(CsrfID)

	err := sr.repo.Del(ctx, CsrfID).Err()
	if err != nil {
		sr.logger.Error(err.Error())
		return err
	}
	return nil
}

func (sr *CsrfRepositoryService) GetCsrf(ctx context.Context, CsrfID string) (string, error) {

	CsrfID = utils.Sanitize(CsrfID)

	email, err := sr.repo.Get(ctx, CsrfID).Result()
	if err != nil {
		sr.logger.Error(err.Error())
		return "", err
	}

	email = utils.Sanitize(email)

	return email, nil
}
