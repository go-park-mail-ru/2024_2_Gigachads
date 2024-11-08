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

func setupRedis(t *testing.T) (*redis.Client, func()) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("ошибка запуска miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, func() {
		client.Close()
		mr.Close()
	}
}

func TestSessionRepositoryService_CreateSession(t *testing.T) {
	t.Run("успешное создание сессии", func(t *testing.T) {
		client, cleanup := setupRedis(t)
		defer cleanup()

		mockLogger := logger.NewLogger()
		repo := NewSessionRepositoryService(client, mockLogger)
		ctx := context.Background()

		email := "test@example.com"

		session, err := repo.CreateSession(ctx, email)
		assert.NoError(t, err)
		assert.NotEmpty(t, session.ID)
		assert.Equal(t, email, session.UserLogin)
		assert.NotZero(t, session.Time)
		assert.Equal(t, "email", session.Name)

		val, err := client.Get(ctx, session.ID).Result()
		assert.NoError(t, err)
		assert.Equal(t, email, val)

		ttl, err := client.TTL(ctx, session.ID).Result()
		assert.NoError(t, err)
		assert.True(t, ttl > 0 && ttl <= 24*time.Hour)
	})

	t.Run("ошибка при создании сессии с пустым email", func(t *testing.T) {
		client, cleanup := setupRedis(t)
		defer cleanup()

		mockLogger := logger.NewLogger()
		repo := NewSessionRepositoryService(client, mockLogger)
		ctx := context.Background()

		email := ""

		session, err := repo.CreateSession(ctx, email)
		assert.NoError(t, err)
		assert.NotEmpty(t, session.ID)

		val, err := client.Get(ctx, session.ID).Result()
		assert.NoError(t, err)
		assert.Empty(t, val)
	})

	t.Run("ошибка Redis при создании сессии", func(t *testing.T) {
		client, cleanup := setupRedis(t)
		mockLogger := logger.NewLogger()
		repo := NewSessionRepositoryService(client, mockLogger)
		ctx := context.Background()

		cleanup()

		email := "test@example.com"
		session, err := repo.CreateSession(ctx, email)
		assert.Error(t, err)
		assert.Nil(t, session)
	})

	t.Run("создание сессии с специальными символами", func(t *testing.T) {
		client, cleanup := setupRedis(t)
		defer cleanup()

		mockLogger := logger.NewLogger()
		repo := NewSessionRepositoryService(client, mockLogger)
		ctx := context.Background()

		email := "test+special@example.com<script>"

		session, err := repo.CreateSession(ctx, email)
		assert.NoError(t, err)
		assert.NotEmpty(t, session.ID)

		val, err := client.Get(ctx, session.ID).Result()
		assert.NoError(t, err)
		assert.NotContains(t, val, "<script>")
	})
}

func TestSessionRepositoryService_DeleteSession(t *testing.T) {
	t.Run("успешное удаление сессии", func(t *testing.T) {
		client, cleanup := setupRedis(t)
		defer cleanup()

		mockLogger := logger.NewLogger()
		repo := NewSessionRepositoryService(client, mockLogger)
		ctx := context.Background()

		sessionID := "test-session"
		err := client.Set(ctx, sessionID, "test@example.com", time.Hour).Err()
		assert.NoError(t, err)

		err = repo.DeleteSession(ctx, sessionID)
		assert.NoError(t, err)

		exists, err := client.Exists(ctx, sessionID).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), exists)
	})

	t.Run("ошибка Redis при удалении сессии", func(t *testing.T) {
		client, cleanup := setupRedis(t)
		mockLogger := logger.NewLogger()
		repo := NewSessionRepositoryService(client, mockLogger)
		ctx := context.Background()

		cleanup() // Закрываем Redis перед тестом

		sessionID := "test-session"
		err := repo.DeleteSession(ctx, sessionID)
		assert.Error(t, err)
	})
}

func TestSessionRepositoryService_GetSession(t *testing.T) {
	client, cleanup := setupRedis(t)
	defer cleanup()

	mockLogger := logger.NewLogger()
	repo := NewSessionRepositoryService(client, mockLogger)
	ctx := context.Background()

	t.Run("успешное получение сессии", func(t *testing.T) {
		sessionID := "test-session"
		expectedEmail := "test@example.com"

		err := client.Set(ctx, sessionID, expectedEmail, time.Hour).Err()
		assert.NoError(t, err)

		email, err := repo.GetSession(ctx, sessionID)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmail, email)
	})

	t.Run("получение несуществующей сессии", func(t *testing.T) {
		sessionID := "nonexistent-session"

		email, err := repo.GetSession(ctx, sessionID)
		assert.Error(t, err)
		assert.Equal(t, "", email)
	})

	t.Run("получение сессии с некорректным ID", func(t *testing.T) {
		sessionID := ""

		email, err := repo.GetSession(ctx, sessionID)
		assert.Error(t, err)
		assert.Equal(t, "", email)
	})
}

func TestNewSessionRepositoryService(t *testing.T) {
	client, cleanup := setupRedis(t)
	defer cleanup()

	mockLogger := logger.NewLogger()
	repo := NewSessionRepositoryService(client, mockLogger)
	assert.NotNil(t, repo)
}
