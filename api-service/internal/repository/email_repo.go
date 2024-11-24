package repository

import (
	"database/sql"
	"errors"
	"mail/api-service/internal/models"
	"mail/api-service/pkg/logger"
	"mail/api-service/pkg/utils"
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

func (er *EmailRepositoryService) Inbox(email string) ([]models.Email, error) { //больше не используется

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

func (er *EmailRepositoryService) GetSentEmails(senderEmail string) ([]models.Email, error) { //больше не используется

	senderEmail = utils.Sanitize(senderEmail)

	rows, err := er.repo.Query(
		`SELECT t.id, t.sender_email, t.recipient_email, m.title, 
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
		`SELECT id FROM folder WHERE user_id = $1 AND name = "Отправленные"`,
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
		email.Sending_date, email.IsRead, messageID,
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
		er.logger.Error(err.Error())
		return err
	}
	var recipientFolderID int
	err = tx.QueryRow(
		`SELECT id FROM folder WHERE user_id = $1 AND name = "Входящие"`,
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

func (er *EmailRepositoryService) DeleteEmails(userEmail string, messageIDs []int, folder string) error {
	tx, err := er.repo.Begin()
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	defer tx.Rollback()

	var query string
	switch folder {
	case "inbox":
		query = `DELETE FROM email_transaction 
				 WHERE message_id = ANY($1) 
				 AND recipient_email = $2`
	case "sent":
		query = `DELETE FROM email_transaction 
				 WHERE message_id = ANY($1) 
				 AND sender_email = $2`
	default:
		return errors.New("неизвестная папка")
	}

	_, err = tx.Exec(query, pq.Array(messageIDs), userEmail)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	return tx.Commit()
}

func (er *EmailRepositoryService) GetFolders(email string) ([]string, error) {
	email = utils.Sanitize(email)

	rows, err := er.repo.Query(
		`SELECT f.name
		 FROM folder AS f
		 JOIN profile AS p ON f.user_id = p.id
		 WHERE p.email = $1`, email)
	if err != nil {
		er.logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	res := make([]string, 0)
	for rows.Next() {
		var folder string
		err := rows.Scan(&folder)
		if err != nil {
			er.logger.Error(err.Error())
			return nil, err
		}
		folder = utils.Sanitize(folder)
		res = append(res, folder)
	}
	return res, nil
}

func (er *EmailRepositoryService) GetFolderEmails(email string, folderName string) ([]models.Email, error) {

	email = utils.Sanitize(email)
	folderName = utils.Sanitize(folderName)

	rows, err := er.repo.Query(
		`SELECT t.id, t.sender_email, t.recipient_email, m.title, 
		 t.sending_date, t.isread, m.description
		 FROM email_transaction AS t
		 JOIN message AS m ON t.message_id = m.id
		 JOIN folder AS f ON t.folder_id = f.id
		 JOIN profile AS p ON f.user_id = p.id
		 WHERE f.name = $2
		 AND p.email = $1
		 ORDER BY t.sending_date DESC`, email, folderName)
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

func (er *EmailRepositoryService) CreateFolder(email string, folderName string) error {

	email = utils.Sanitize(email)
	folderName = utils.Sanitize(folderName)

	tx, err := er.repo.Begin()
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	defer tx.Rollback()

	var userID int

	err = tx.QueryRow(
		`SELECT id FROM profile WHERE email = $1`,
		email, 
	).Scan(&userID)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO folder (user_id, name) 
		VALUES ($1, $2)`,
		userID, folderName,
	)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	return tx.Commit()
}

func (er *EmailRepositoryService) DeleteFolder(email string, folderName string) error {

	email = utils.Sanitize(email)
	folderName = utils.Sanitize(folderName)

	tx, err := er.repo.Begin()
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	defer tx.Rollback()

	var userID int

	err = tx.QueryRow(
		`SELECT id FROM profile WHERE email = $1`,
		email, 
	).Scan(&userID)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM folder 
		 WHERE name = $2
		 AND user_id = $1`,
		userID, folderName,
	)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	return tx.Commit()
}

func (er *EmailRepositoryService) RenameFolder(email string, folderName string, newFolderName string) error {

	email = utils.Sanitize(email)
	folderName = utils.Sanitize(folderName)
	newFolderName = utils.Sanitize(newFolderName)

	tx, err := er.repo.Begin()
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	defer tx.Rollback()

	var userID int

	err = tx.QueryRow(
		`SELECT id FROM profile WHERE email = $1`,
		email, 
	).Scan(&userID)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	var folderID int

	err = tx.QueryRow(
		`SELECT id FROM folder WHERE user_id = $1 AND name = $2`,
		userID, folderName, 
	).Scan(&folderID)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	_, err = tx.Exec(
		`UPDATE folder
		 SET name = $2
		 WHERE message_id = $1`,
		folderID, newFolderName,
	)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	return tx.Commit()
}

func (er *EmailRepositoryService) ChangeEmailFolder(id int, email string, folderName string) error {
	email = utils.Sanitize(email)
	folderName = utils.Sanitize(folderName)

	tx, err := er.repo.Begin()
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	defer tx.Rollback()

	var userID int

	err = tx.QueryRow(
		`SELECT id FROM profile WHERE email = $1`,
		email, 
	).Scan(&userID)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	var folderID int

	err = tx.QueryRow(
		`SELECT id FROM folder WHERE user_id = $1 AND name = $2`,
		userID, folderName, 
	).Scan(&folderID)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	_, err = tx.Exec(
		`UPDATE email_transaction
			SET folder_id = $2
			WHERE message_id = $1`,
		id, folderID,
	)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	return tx.Commit()
}

func (er *EmailRepositoryService) CreateDraft(email models.Email) error {
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
		`SELECT id FROM folder WHERE user_id = $1 AND name = "Черновики"`,
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
		email.Sending_date, email.IsRead, messageID,
		parentID, senderFolderID,
	)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	return tx.Commit()
}

func (er *EmailRepositoryService) UpdateDraft(email models.Draft) error {
	
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
		`SELECT message_id FROM email_transaction WHERE id = $1`,
		email.ID, 
	).Scan(&messageID)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}

	_, err = tx.Exec(
		`UPDATE message
			SET title = $2, description = $3
			WHERE mid = $1`,
		messageID, email.Title, email.Description,
	)
	if err != nil {
		er.logger.Error(err.Error())
		return err
	}
	return tx.Commit()
}
