package models

import (
	"time"
)

type Session struct {
	Name      string
	ID        string
	Time      time.Time
	UserLogin string
}

type SessionRepository interface {
	CreateSession(mail string) (*Session, error)
	DeleteSession(sessionID string) error
	GetSession(sessionID string) (*Session, error)
}
