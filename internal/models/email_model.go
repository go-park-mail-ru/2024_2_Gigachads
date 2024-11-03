package models

import "time"

type Email struct {
	ID           int       `json:"id"`
	ParentID     int       `json:"parentID"`
	Sender_email string    `json:"author"`
	Recipient    string    `json:"recipient"`
	Title        string    `json:"title"`
	IsRead       bool      `json:"is_read"`
	Sending_date time.Time `json:"date"`
	Description  string    `json:"description"`
}

type EmailUseCase interface {
	Inbox(id string) ([]Email, error)
	SendEmail(from string, to []string, subject string, body string) error
	ForwardEmail(from string, to []string, originalEmail Email) error
	ReplyEmail(from string, to string, originalEmail Email, replyText string) error
	GetEmailByID(id int) (Email, error)
	FetchEmailsViaPOP3() error
	ChangeStatus(id int, status string) error
	GetSentEmails(senderEmail string) ([]Email, error)
	SaveEmail(email Email) error
}

type EmailRepository interface {
	Inbox(id string) ([]Email, error)
	GetEmailByID(id int) (Email, error)
	SaveEmail(email Email) error
	ChangeStatus(id int, status string) error
	GetSentEmails(senderEmail string) ([]Email, error)
}

type SMTPRepository interface {
	SendEmail(from string, to []string, subject string, body string) error
}

type POP3Repository interface {
	Connect() error
	FetchEmails(EmailRepository) error
	Quit() error
}
