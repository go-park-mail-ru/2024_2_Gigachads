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
		`SELECT t.id, t.sender_email, t.recipient_email, m.title, 
		 t.sending_date, t.isread, m.description
		 FROM email_transaction AS t
		 JOIN message AS m ON t.message_id = m.id
		 WHERE t.recipient_email = $1
		 ORDER BY t.sending_date DESC`, uEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]models.Email, 0)
	for rows.Next() {
		email := models.Email{}
		err := rows.Scan(
			&email.ID,
			&email.Sender_email,
			&email.Recipient,
			&email.Title,
			&email.Sending_date,
			&email.IsRead,
			&email.Description,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, email)
	}
	return res, nil
}

func (er *EmailRepositoryService) GetSentEmails(senderEmail string) ([]models.Email, error) {
	rows, err := er.repo.Query(
		`SELECT t.id, t.sender_email, t.recipient_email, m.title, 
		 t.sending_date, t.isread, m.description
		 FROM email_transaction AS t
		 JOIN message AS m ON t.message_id = m.id
		 WHERE t.sender_email = $1
		 ORDER BY t.sending_date DESC`, senderEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]models.Email, 0)
	for rows.Next() {
		email := models.Email{}
		err := rows.Scan(
			&email.ID,
			&email.Sender_email,
			&email.Recipient,
			&email.Title,
			&email.Sending_date,
			&email.IsRead,
			&email.Description,
		)
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

	query := `
	SELECT t.id, t.sender_email, t.recipient_email, m.title, 
	t.isread, t.sending_date, m.description
	FROM email_transaction AS t
	JOIN message AS m ON t.message_id = m.id
	WHERE t.id = $1
	`

	var email models.Email
	err := er.repo.QueryRow(query, id).
		Scan(&email.ID, &email.Sender_email, &email.Recipient,
			&email.Title, &email.IsRead, &email.Sending_date,
			&email.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Email{}, errors.New("email not found")
		}
		return models.Email{}, err
	}
	return email, nil
}

func (er *EmailRepositoryService) SaveEmail(email models.Email) error {
	er.mu.Lock()
	defer er.mu.Unlock()

	tx, err := er.repo.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var messageID int
	err = tx.QueryRow(
		`INSERT INTO message (title, description) 
		VALUES ($1, $2) RETURNING id`,
		email.Title, email.Description,
	).Scan(&messageID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO email_transaction 
		(sender_email, recipient_email, title, sending_date, isread, message_id)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		email.Sender_email, email.Recipient, email.Title,
		email.Sending_date, email.IsRead, messageID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (er *EmailRepositoryService) ChangeStatus(id int, status string) error {
	er.mu.Lock()
	defer er.mu.Unlock()

	tx, err := er.repo.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if status == "true" {
		_, err = tx.Exec(
			`UPDATE email_transaction
			SET isread = TRUE
			WHERE id = $1`,
			id,
		)
	} else {
		_, err = tx.Exec(
			`UPDATE email_transaction
			SET isread = FALSE
			WHERE id = $1`,
			id,
		)
	}

	if err != nil {
		return err
	}
	return tx.Commit()
}
