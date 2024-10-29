package repository

import (
	"database/sql"
	models "mail/internal/models"
)

type EmailRepositoryService struct {
	repo *sql.DB
}

func NewEmailRepositoryService(db *sql.DB) models.EmailRepository {
	return &EmailRepositoryService{repo: db}
}

func (er *EmailRepositoryService) Inbox(uEmail string) ([]models.Email, error) {
	rows, err := er.repo.Query(
		`SELECT email_transaction.sender_email, email_transaction.sending_date, email_transaction.isread, message.title, message.decription FROM email_transaction AS t
		JOIN message AS m ON t.message_id = m.id
		WHERE t.recipient_email = '$1'`, uEmail)
	if err != nil {
		return nil, err
	}
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
