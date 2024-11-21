package models

import (
	"context"
	"time"
)

type Csrf struct {
	Name      string
	ID        string
	Time      time.Time
	UserLogin string
}

type CsrfRepository interface {
	CreateCsrf(ctx context.Context, mail string) (*Csrf, error)
	DeleteCsrf(ctx context.Context, CsrfID string) error
	GetCsrf(ctx context.Context, CsrfID string) (string, error)
}