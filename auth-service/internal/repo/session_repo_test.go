package repository

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...interface{})  {}
func (m *mockLogger) Error(msg string, args ...interface{}) {}
func (m *mockLogger) Fatal(msg string, args ...interface{}) {}
func (m *mockLogger) Debug(msg string, args ...interface{}) {}
func (m *mockLogger) Warn(msg string, args ...interface{})  {}

func setupRedis(t *testing.T) (*redis.Client, func()) {
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

func TestSessionRepository_CreateSession(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupRedis(t)
			defer cleanup()

			repo := &SessionRepositoryService{
				repo:   client,
				logger: logger,
			}

			session, err := repo.CreateSession(ctx, tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, session)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, session)
				assert.Equal(t, tt.email, session.UserLogin)
				assert.NotEmpty(t, session.ID)
				assert.True(t, session.Time.After(time.Now().Add(-time.Second)))
			}
		})
	}
}

func TestSessionRepository_GetSession(t *testing.T) {
	ctx := context.Background()
	logger := &mockLogger{}

	tests := []struct {
		name      string
		sessionID string
		setupFn   func(*redis.Client)
		wantEmail string
		wantErr   bool
	}{
		{
			name:      "Success",
			sessionID: "test-session-id",
			setupFn: func(client *redis.Client) {
				client.Set(ctx, "test-session-id", "test@example.com", time.Hour)
			},
			wantEmail: "test@example.com",
			wantErr:   false,
		},
		{
			name:      "Not found",
			sessionID: "nonexistent-id",
			setupFn:   func(client *redis.Client) {},
			wantEmail: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupRedis(t)
			defer cleanup()

			tt.setupFn(client)

			repo := &SessionRepositoryService{
				repo:   client,
				logger: logger,
			}

			email, err := repo.GetSession(ctx, tt.sessionID)
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

func TestSessionRepository_DeleteSession(t *testing.T) {
	ctx := context.Background()
	logger := &mockLogger{}

	tests := []struct {
		name      string
		sessionID string
		setupFn   func(*redis.Client)
		wantErr   bool
	}{
		{
			name:      "Success",
			sessionID: "test-session-id",
			setupFn: func(client *redis.Client) {
				client.Set(ctx, "test-session-id", "test@example.com", time.Hour)
			},
			wantErr: false,
		},
		{
			name:      "Not found",
			sessionID: "nonexistent-id",
			setupFn:   func(client *redis.Client) {},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupRedis(t)
			defer cleanup()

			tt.setupFn(client)

			repo := &SessionRepositoryService{
				repo:   client,
				logger: logger,
			}

			err := repo.DeleteSession(ctx, tt.sessionID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
