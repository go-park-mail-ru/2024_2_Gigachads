package repository

import (
	"database/sql"
	"errors"
	"mail/internal/models"
	"sync"
)

type EmailRepositoryService struct {
	repo *sql.DB
	mu   sync.RWMutex
}

func NewEmailRepositoryService(db *sql.DB) *EmailRepositoryService {
	return &EmailRepositoryService{repo: db}
}

func (er *EmailRepositoryService) Inbox(uEmail string) ([]models.Email, error) {
	rows, err := er.repo.Query(
		`SELECT email_transaction.sender_email, email_transaction.sending_date, email_transaction.isread, message.title, message.description 
		FROM email_transaction AS t
		JOIN message AS m ON t.message_id = m.id
		WHERE t.recipient_email = $1`, uEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]models.Email, 0)
	for rows.Next() {
		email := models.Email{}
		err := rows.Scan(&email.Sender_email, &email.Sending_date, &email.IsRead, &email.Title, &email.Description)
		if err != nil {
			return nil, err
		}
		res = append(res, email)
	}
	return res, nil
}

func (er *EmailRepositoryService) GetEmailByID(id int) (models.Email, error) {
	er.mu.RLock()
	defer er.mu.RUnlock()

	var email models.Email
	err := er.repo.QueryRow(
		`SELECT id, sender_email, recipient, title, is_read, sending_date, description 
		FROM emails WHERE id = $1`, id).
		Scan(&email.ID, &email.Sender_email, &email.Recipient, &email.Title, &email.IsRead, &email.Sending_date, &email.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Email{}, errors.New("email not found")
		}
		return models.Email{}, err
	}
	return email, nil
}
