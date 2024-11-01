package models

import (
	"context"
	"time"
)

type Email struct {
	Sender_email string    `json:"author"`
	Title        string    `json:"title"`
	IsRead       bool      `json:"is_read"`
	Sending_date time.Time `json:"date"`
	Description  string    `json:"description"`
}

type EmailUseCase interface {
	Inbox(ctx context.Context, id string) ([]Email, error)
}

type EmailRepository interface {
	Inbox(id string) ([]Email, error)
}
