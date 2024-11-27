package repository

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRedisForCsrf(t *testing.T) (*redis.Client, func()) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, func() {
		client.Close()
		mr.Close()
	}
}

func TestCsrfRepository_CreateCsrf(t *testing.T) {
	ctx := context.Background()
	logger := &mockLogger{}

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Success",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Empty email",
			email:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupRedisForCsrf(t)
			defer cleanup()

			repo := &CsrfRepositoryService{
				repo:   client,
				logger: logger,
			}

			csrf, err := repo.CreateCsrf(ctx, tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, csrf)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, csrf)
				assert.Equal(t, "csrf", csrf.Name)
				assert.Equal(t, tt.email, csrf.UserLogin)
				assert.NotEmpty(t, csrf.ID)
				assert.True(t, csrf.Time.After(time.Now()))

				// Проверяем, что значение сохранено в Redis
				val, err := client.Get(ctx, csrf.ID).Result()
				assert.NoError(t, err)
				assert.Equal(t, tt.email, val)
			}
		})
	}
}

func TestCsrfRepository_GetCsrf(t *testing.T) {
	ctx := context.Background()
	logger := &mockLogger{}

	tests := []struct {
		name      string
		csrfID    string
		setupFn   func(*redis.Client)
		wantEmail string
		wantErr   bool
	}{
		{
			name:   "Success",
			csrfID: "test-csrf-id",
			setupFn: func(client *redis.Client) {
				client.Set(ctx, "test-csrf-id", "test@example.com", time.Hour)
			},
			wantEmail: "test@example.com",
			wantErr:   false,
		},
		{
			name:      "Not found",
			csrfID:    "nonexistent-id",
			setupFn:   func(client *redis.Client) {},
			wantEmail: "",
			wantErr:   true,
		},
		{
			name:   "Empty email in Redis",
			csrfID: "empty-email-id",
			setupFn: func(client *redis.Client) {
				client.Set(ctx, "empty-email-id", "", time.Hour)
			},
			wantEmail: "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupRedisForCsrf(t)
			defer cleanup()

			tt.setupFn(client)

			repo := &CsrfRepositoryService{
				repo:   client,
				logger: logger,
			}

			email, err := repo.GetCsrf(ctx, tt.csrfID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, email)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantEmail, email)
			}
		})
	}
}

func TestCsrfRepository_DeleteCsrf(t *testing.T) {
	ctx := context.Background()
	logger := &mockLogger{}

	tests := []struct {
		name    string
		csrfID  string
		setupFn func(*redis.Client)
		wantErr bool
	}{
		{
			name:   "Success",
			csrfID: "test-csrf-id",
			setupFn: func(client *redis.Client) {
				client.Set(ctx, "test-csrf-id", "test@example.com", time.Hour)
			},
			wantErr: false,
		},
		{
			name:    "Not found",
			csrfID:  "nonexistent-id",
			setupFn: func(client *redis.Client) {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupRedisForCsrf(t)
			defer cleanup()

			tt.setupFn(client)

			repo := &CsrfRepositoryService{
				repo:   client,
				logger: logger,
			}

			err := repo.DeleteCsrf(ctx, tt.csrfID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Проверяем, что значение удалено из Redis
				exists, err := client.Exists(ctx, tt.csrfID).Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), exists)
			}
		})
	}
}

func TestNewCsrfRepositoryService(t *testing.T) {
	client, cleanup := setupRedisForCsrf(t)
	defer cleanup()
	logger := &mockLogger{}

	repo := NewCsrfRepositoryService(client, logger)
	require.NotNil(t, repo)
}
