package repo

import (
	"database/sql"
	"mail/api-service/pkg/logger"
	"mail/api-service/pkg/utils"
	"mail/smtp-service/internal/models"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type EmailRepositoryService struct {
	repo   *sql.DB
	redis  *redis.Client
	logger logger.Logable
}

func NewEmailRepositoryService(db *sql.DB, r *redis.Client, l logger.Logable) *EmailRepositoryService {
	return &EmailRepositoryService{repo: db, redis: r, logger: l}
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

	var parentID interface{}
	if email.ParentID == 0 {
		parentID = nil
	} else {
		parentID = email.ParentID
	}
	var senderID int
	err = tx.QueryRow(
		`SELECT id FROM profile WHERE email = $1`,
		email.Sender_email,
	).Scan(&senderID)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	var senderFolderID int
	err = tx.QueryRow(
		`SELECT id FROM folder WHERE user_id = $1 AND name = 'Отправленные'`,
		senderID,
	).Scan(&senderFolderID)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	_, err = tx.Exec(
		`INSERT INTO email_transaction 
        (sender_email, recipient_email, sending_date, isread, message_id, parent_transaction_id, folder_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		email.Sender_email, email.Recipient,
		email.Sending_date, true /*email.IsRead*/, messageID,
		parentID, senderFolderID,
	)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	var recipientID int
	err = tx.QueryRow(
		`SELECT id FROM profile WHERE email = $1`,
		email.Recipient,
	).Scan(&recipientID)
	if err != nil {
		if err != sql.ErrNoRows {
			er.logger.Error(err.Error())
			return err
		}
	} else {
		var recipientFolderID int
		err = tx.QueryRow(
			`SELECT id FROM folder WHERE user_id = $1 AND name = 'Входящие'`,
			recipientID,
		).Scan(&recipientFolderID)
		if err != nil {
			er.logger.Error(err.Error())
			return err
		}

		_, err = tx.Exec(
			`INSERT INTO email_transaction 
            (sender_email, recipient_email, sending_date, isread, message_id, parent_transaction_id, folder_id)
            VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			email.Sender_email, email.Recipient,
			email.Sending_date, email.IsRead, messageID,
			parentID, recipientFolderID,
		)
		if err != nil {
			er.logger.Error(err.Error())
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	return nil
}
