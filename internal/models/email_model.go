package models

import "time"

type Email struct {
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Text        string    `json:"text"`
	Badge_text  string    `json:"badge_text"`
	Badge_type  string    `json:"badge_type"`
	Date        time.Time `json:"date"`
}

type EmailUseCase interface {
	Inbox(id string) ([]Email, error)
}

type EmailRepository interface {
	Inbox(id string) ([]Email, error)
}
