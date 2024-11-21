package repository

import (
	"context"
	"mail/pkg/logger"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupRedisForCsrf(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("ошибка запуска miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

func TestCsrfRepositoryService_CreateCsrf(t *testing.T) {
	t.Run("успешное создание CSRF", func(t *testing.T) {
		client, mr := setupRedisForCsrf(t)
		defer func() {
			client.Close()
			mr.Close()
		}()

		mockLogger := logger.NewLogger()
		repo := NewCsrfRepositoryService(client, mockLogger)
		ctx := context.Background()

		email := "test@example.com"

		csrf, err := repo.CreateCsrf(ctx, email)
		assert.NoError(t, err)
		assert.NotEmpty(t, csrf.ID)
		assert.Equal(t, email, csrf.UserLogin)
		assert.NotZero(t, csrf.Time)
		assert.Equal(t, "csrf", csrf.Name)

		val, err := client.Get(ctx, csrf.ID).Result()
		assert.NoError(t, err)
		assert.Equal(t, email, val)
	})

	t.Run("ошибка Redis при создании CSRF", func(t *testing.T) {
		client, mr := setupRedisForCsrf(t)
		mockLogger := logger.NewLogger()
		repo := NewCsrfRepositoryService(client, mockLogger)
		ctx := context.Background()

		mr.Close()
		client.Close()

		email := "test@example.com"
		csrf, err := repo.CreateCsrf(ctx, email)
		assert.Error(t, err)
		assert.Nil(t, csrf)
	})
}

func TestCsrfRepositoryService_DeleteCsrf(t *testing.T) {
	t.Run("успешное удаление CSRF", func(t *testing.T) {
		client, mr := setupRedisForCsrf(t)
		defer func() {
			client.Close()
			mr.Close()
		}()

		mockLogger := logger.NewLogger()
		repo := NewCsrfRepositoryService(client, mockLogger)
		ctx := context.Background()

		csrfID := "test-csrf"
		err := client.Set(ctx, csrfID, "test@example.com", time.Hour).Err()
		assert.NoError(t, err)

		err = repo.DeleteCsrf(ctx, csrfID)
		assert.NoError(t, err)

		exists, err := client.Exists(ctx, csrfID).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), exists)
	})

	t.Run("ошибка Redis при удалении CSRF", func(t *testing.T) {
		client, mr := setupRedisForCsrf(t)
		mockLogger := logger.NewLogger()
		repo := NewCsrfRepositoryService(client, mockLogger)
		ctx := context.Background()

		mr.Close()
		client.Close()

		csrfID := "test-csrf"
		err := repo.DeleteCsrf(ctx, csrfID)
		assert.Error(t, err)
	})

	t.Run("удаление несуществующего CSRF", func(t *testing.T) {
		client, mr := setupRedisForCsrf(t)
		defer func() {
			client.Close()
			mr.Close()
		}()

		mockLogger := logger.NewLogger()
		repo := NewCsrfRepositoryService(client, mockLogger)
		ctx := context.Background()

		csrfID := "nonexistent-csrf"
		err := repo.DeleteCsrf(ctx, csrfID)
		assert.NoError(t, err)
	})
}

func TestCsrfRepositoryService_GetCsrf(t *testing.T) {
	t.Run("успешное получение CSRF", func(t *testing.T) {
		client, mr := setupRedisForCsrf(t)
		defer func() {
			client.Close()
			mr.Close()
		}()

		mockLogger := logger.NewLogger()
		repo := NewCsrfRepositoryService(client, mockLogger)
		ctx := context.Background()

		csrfID := "test-csrf"
		expectedEmail := "test@example.com"

		err := client.Set(ctx, csrfID, expectedEmail, time.Hour).Err()
		assert.NoError(t, err)

		email, err := repo.GetCsrf(ctx, csrfID)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmail, email)
	})

	t.Run("получение несуществующего CSRF", func(t *testing.T) {
		client, mr := setupRedisForCsrf(t)
		defer func() {
			client.Close()
			mr.Close()
		}()

		mockLogger := logger.NewLogger()
		repo := NewCsrfRepositoryService(client, mockLogger)
		ctx := context.Background()

		csrfID := "nonexistent-csrf"

		email, err := repo.GetCsrf(ctx, csrfID)
		assert.Error(t, err)
		assert.Empty(t, email)
	})

	t.Run("ошибка Redis при получении CSRF", func(t *testing.T) {
		client, mr := setupRedisForCsrf(t)
		mockLogger := logger.NewLogger()
		repo := NewCsrfRepositoryService(client, mockLogger)
		ctx := context.Background()

		mr.Close()
		client.Close()

		csrfID := "test-csrf"
		email, err := repo.GetCsrf(ctx, csrfID)
		assert.Error(t, err)
		assert.Empty(t, email)
	})
}

func TestNewCsrfRepositoryService(t *testing.T) {
	client, mr := setupRedisForCsrf(t)
	defer func() {
		client.Close()
		mr.Close()
	}()

	mockLogger := logger.NewLogger()
	repo := NewCsrfRepositoryService(client, mockLogger)
	assert.NotNil(t, repo)
}
