package models

import (
	"time"
)

type EmailRepositorySMTP interface {
	SaveEmail(email Email) error
}

type SMTPRepository interface {
	SendEmail(from string, to []string, subject string, body string) error
}

type POP3Repository interface {
	Connect() error
	FetchEmails(EmailRepositorySMTP) error
	Quit() error
}

type Email struct {
	ID           int       `json:"id"`
	ParentID     int       `json:"parentID"`
	Sender_email string    `json:"sender"`
	Recipient    string    `json:"recipient"`
	Title        string    `json:"title"`
	IsRead       bool      `json:"isRead"`
	Sending_date time.Time `json:"date"`
	Description  string    `json:"description"`
}
