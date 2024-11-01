package models

import (
	"context"
	"time"
)

type Session struct {
	Name      string
	ID        string
	Time      time.Time
	UserLogin string
}

type SessionRepository interface {
	CreateSession(ctx context.Context, mail string) (*Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
	GetSession(ctx context.Context, sessionID string) (string, error)
}
