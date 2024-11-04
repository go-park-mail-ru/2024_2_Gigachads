package repository

import (
	"database/sql"
	"errors"
	"mail/internal/models"
	"mail/pkg/logger"
	"mail/pkg/utils"
	"strconv"
	"sync"

	"github.com/lib/pq"
)

type EmailRepositoryService struct {
	repo   *sql.DB
	mu     sync.RWMutex
	logger logger.Logable
}

func NewEmailRepositoryService(db *sql.DB, l logger.Logable) *EmailRepositoryService {
	return &EmailRepositoryService{repo: db, logger: l}
}

func (er *EmailRepositoryService) Inbox(email string) ([]models.Email, error) {

	email = utils.Sanitize(email)

	rows, err := er.repo.Query(
		`SELECT t.id, t.sender_email, t.recipient_email, m.title, 
		 t.sending_date, t.isread, m.description
		 FROM email_transaction AS t
		 JOIN message AS m ON t.message_id = m.id
		 WHERE t.recipient_email = $1
		 ORDER BY t.sending_date DESC`, email)
	if err != nil {
		er.logger.Error(err.Error())
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
		email.Sender_email = utils.Sanitize(email.Sender_email)
		email.Recipient = utils.Sanitize(email.Recipient)
		email.Title = utils.Sanitize(email.Title)
		email.Description = utils.Sanitize(email.Description)
		if err != nil {
			er.logger.Error(err.Error())
			return nil, err
		}
		res = append(res, email)
	}
	return res, nil
}

func (er *EmailRepositoryService) GetSentEmails(senderEmail string) ([]models.Email, error) {

	senderEmail = utils.Sanitize(senderEmail)

	rows, err := er.repo.Query(
		`SELECT t.id, t.sender_email, m.title, 
		 t.sending_date, t.isread, m.description
		 FROM email_transaction AS t
		 JOIN message AS m ON t.message_id = m.id
		 WHERE t.sender_email = $1
		 ORDER BY t.sending_date DESC`, senderEmail)
	if err != nil {
		er.logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	res := make([]models.Email, 0)
	for rows.Next() {
		email := models.Email{}
		err := rows.Scan(
			&email.ID,
			&email.Sender_email,
			&email.Title,
			&email.Sending_date,
			&email.IsRead,
			&email.Description,
		)
		email.Sender_email = utils.Sanitize(email.Sender_email)
		email.Recipient = utils.Sanitize(email.Recipient)
		email.Title = utils.Sanitize(email.Title)
		email.Description = utils.Sanitize(email.Description)
		if err != nil {
			er.logger.Error(err.Error())
			return nil, err
		}
		res = append(res, email)
	}
	return res, nil
}

func (er *EmailRepositoryService) GetEmailByID(id int) (models.Email, error) {
	query := `
	SELECT t.id, parent_transaction_id, t.sender_email, t.recipient_email, m.title, 
	t.isread, t.sending_date, m.description
	FROM email_transaction AS t
	JOIN message AS m ON t.message_id = m.id
	WHERE m.id = $1
	`

	var email models.Email
	var parentIdNullString sql.NullString
	err := er.repo.QueryRow(query, id).
		Scan(&email.ID, &parentIdNullString, &email.Sender_email, &email.Recipient,
			&email.Title, &email.IsRead, &email.Sending_date,
			&email.Description)
	if err != nil {
		er.logger.Error(err.Error())
		return models.Email{}, err
	}

	if parentIdNullString.String == "" {
		email.ParentID = 0
	} else {
		email.ParentID, err = strconv.Atoi(parentIdNullString.String)
		if err != nil {
			er.logger.Error(err.Error())
			return models.Email{}, err
		}
	}

	email.Sender_email = utils.Sanitize(email.Sender_email)
	email.Recipient = utils.Sanitize(email.Recipient)
	email.Title = utils.Sanitize(email.Title)
	email.Description = utils.Sanitize(email.Description)

	if err != nil {
		er.logger.Error(err.Error())
		if err == sql.ErrNoRows {
			return models.Email{}, errors.New("email not found")
		}
		return models.Email{}, err
	}
	return email, nil
}

func (er *EmailRepositoryService) SaveEmail(email models.Email) error {

	email.Sender_email = utils.Sanitize(email.Sender_email)
	email.Recipient = utils.Sanitize(email.Recipient)
	email.Title = utils.Sanitize(email.Title)
	email.Description = utils.Sanitize(email.Description)

	tx, err := er.repo.Begin()
	if err != nil {
		er.logger.Error(err.Error())
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
		er.logger.Error(err.Error())
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO email_transaction 
		(sender_email, recipient_email, sending_date, isread, message_id)
		VALUES ($1, $2, $3, $4, $5)`,
		email.Sender_email, email.Recipient,
		email.Sending_date, email.IsRead, messageID,
	)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	return tx.Commit()
}

func (er *EmailRepositoryService) ChangeStatus(id int, status bool) error {
	tx, err := er.repo.Begin()
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	defer tx.Rollback()

	if status {
		_, err = tx.Exec(
			`UPDATE email_transaction
			SET isread = TRUE
			WHERE message_id = $1`,
			id,
		)
	} else {
		_, err = tx.Exec(
			`UPDATE email_transaction
			SET isread = FALSE
			WHERE message_id = $1`,
			id,
		)
	}

	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	return tx.Commit()
}

func (er *EmailRepositoryService) DeleteEmails(ids []int) error {
	tx, err := er.repo.Begin()
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM email_transaction WHERE id = ANY($1)`

	_, err = tx.Exec(query, pq.Array(ids))
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	return tx.Commit()
}
