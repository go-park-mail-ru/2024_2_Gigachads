package models

import "time"

type Email struct {
	Sender_email string    `json:"author"`
	Title        string    `json:"title"`
	IsRead       bool      `json:"is_read"`
	Sending_date time.Time `json:"date"`
	Description  string    `json:"description"`
}

type EmailUseCase interface {
	Inbox(id string) ([]Email, error)
	SendEmail(from string, to []string, subject string, body string) error
}

type EmailRepository interface {
	Inbox(id string) ([]Email, error)
}

type SMTPRepository interface {
	SendEmail(from string, to []string, subject string, body string) error
}
