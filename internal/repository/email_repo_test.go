package repository

import (
	"database/sql"
	"errors"
	"mail/internal/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)

	email := models.Email{
		Sender_email: "sender@example.com",
		Recipient:    "recipient@example.com",
		Title:        "Тестовое письмо",
		Description:  "Это тело тестового письма.",
		Sending_date: time.Now(),
		IsRead:       false,
	}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO message (title, description) VALUES ($1, $2) RETURNING id`)).
		WithArgs(email.Title, email.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO email_transaction (sender_email, recipient_email, title, sending_date, isread, message_id) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(email.Sender_email, email.Recipient, email.Title, email.Sending_date, email.IsRead, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.SaveEmail(email)
	if err != nil {
		t.Errorf("Неожиданная ошибка: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания были выполнены: %v", err)
	}
}

func TestSaveEmail_InsertMessageError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)

	email := models.Email{
		Sender_email: "sender@example.com",
		Recipient:    "recipient@example.com",
		Title:        "Тестовое письмо",
		Description:  "Это тело тестового письма.",
		Sending_date: time.Now(),
		IsRead:       false,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO message (title, description) VALUES ($1, $2) RETURNING id`)).
		WithArgs(email.Title, email.Description).
		WillReturnError(errors.New("ошибка вставки message"))
	mock.ExpectRollback()

	err = repo.SaveEmail(email)
	if err == nil {
		t.Errorf("Ожидалась ошибка, но её не было")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания были выполнены: %v", err)
	}
}

func TestSaveEmail_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)

	email := models.Email{
		Sender_email: "sender@example.com",
		Recipient:    "recipient@example.com",
		Title:        "Тестовое письмо",
		Description:  "Это тело тестового письма.",
		Sending_date: time.Now(),
		IsRead:       false,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO message (title, description) VALUES ($1, $2) RETURNING id`)).
		WithArgs(email.Title, email.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO email_transaction (sender_email, recipient_email, title, sending_date, isread, message_id) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(email.Sender_email, email.Recipient, email.Title, email.Sending_date, email.IsRead, 1).
		WillReturnError(errors.New("ошибка выполнения Exec"))
	mock.ExpectRollback()

	err = repo.SaveEmail(email)
	if err == nil {
		t.Errorf("Ожидалась ошибка выполнения Exec, но её не было")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания были выполнены: %v", err)
	}
}

func TestGetEmailByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)

	emailID := 1

	query := "SELECT t.id, t.sender_email, t.recipient_email, m.title, t.isread, t.sending_date, m.description FROM email_transaction AS t JOIN message AS m ON t.message_id = m.id WHERE t.id = $1"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(emailID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sender_email", "recipient_email", "title", "isread", "sending_date", "description"}).
			AddRow(emailID, "sender@example.com", "recipient@example.com", "Тестовое письмо", false, time.Now(), "Описание письма"))

	email, err := repo.GetEmailByID(emailID)
	if err != nil {
		t.Errorf("Неожиданная ошибка: %v", err)
	}

	if email.ID != emailID {
		t.Errorf("Неверный ID: получили %d, ожидали %d", email.ID, emailID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания были выполнены: %v", err)
	}
}

func TestGetEmailByID_InvalidID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)

	invalidID := -1

	query := "SELECT t.id, t.sender_email, t.recipient_email, m.title, t.isread, t.sending_date, m.description FROM email_transaction AS t JOIN message AS m ON t.message_id = m.id WHERE t.id = $1"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(invalidID).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetEmailByID(invalidID)
	if err == nil {
		t.Errorf("Ожидалась ошибка для некорректного ID, но её не было")
	}

	if err.Error() != "email not found" {
		t.Errorf("Ожидалась ошибка 'email not found', получили %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания были выполнены: %v", err)
	}
}

func TestInbox(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)

	userEmail := "user@example.com"

	query := `SELECT email_transaction.sender_email, email_transaction.sending_date, email_transaction.isread, message.title, message.description 
	FROM email_transaction AS t
	JOIN message AS m ON t.message_id = m.id
	WHERE t.recipient_email = $1`

	rows := sqlmock.NewRows([]string{"sender_email", "sending_date", "isread", "title", "description"}).
		AddRow("sender1@example.com", time.Now(), false, "Письмо 1", "Описание письма 1").
		AddRow("sender2@example.com", time.Now(), true, "Письмо 2", "Описание письма 2")

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userEmail).
		WillReturnRows(rows)

	emails, err := repo.Inbox(userEmail)
	if err != nil {
		t.Errorf("Неожиданная ошибка: %v", err)
	}

	if len(emails) != 2 {
		t.Errorf("Ожидалось 2 письма, получено %d", len(emails))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания были выполнены: %v", err)
	}
}

func TestInbox_NoEmails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)

	userEmail := "empty@example.com"

	query := `SELECT email_transaction.sender_email, email_transaction.sending_date, email_transaction.isread, message.title, message.description 
	FROM email_transaction AS t
	JOIN message AS m ON t.message_id = m.id
	WHERE t.recipient_email = $1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userEmail).
		WillReturnRows(sqlmock.NewRows([]string{"sender_email", "sending_date", "isread", "title", "description"}))

	emails, err := repo.Inbox(userEmail)
	if err != nil {
		t.Errorf("Неожиданная ошибка: %v", err)
	}

	if len(emails) != 0 {
		t.Errorf("Ожидалось 0 писем, получено %d", len(emails))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания были выполнены: %v", err)
	}
}

func TestInbox_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)

	userEmail := "user@example.com"

	query := `SELECT email_transaction.sender_email, email_transaction.sending_date, email_transaction.isread, message.title, message.description 
	FROM email_transaction AS t
	JOIN message AS m ON t.message_id = m.id
	WHERE t.recipient_email = $1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userEmail).
		WillReturnError(errors.New("ошибка запроса"))

	_, err = repo.Inbox(userEmail)
	if err == nil {
		t.Errorf("Ожидалась ошибка, но её не было")
	}

	if err.Error() != "ошибка запроса" {
		t.Errorf("Ожидалась ошибка 'ошибка запроса', получили %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания были выполнены: %v", err)
	}
}

func TestSaveEmail_BeginTxError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)
	mock.ExpectBegin().WillReturnError(errors.New("ошибка начала транзакции"))

	email := models.Email{
		Sender_email: "sender@example.com",
		Recipient:    "recipient@example.com",
		Title:        "Test",
		Description:  "Test body",
	}

	err = repo.SaveEmail(email)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка начала транзакции")
}

func TestSaveEmail_CommitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)
	email := models.Email{
		Sender_email: "sender@example.com",
		Recipient:    "recipient@example.com",
		Title:        "Test",
		Description:  "Test body",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO message`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO email_transaction`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(errors.New("ошибка коммита"))

	err = repo.SaveEmail(email)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка коммита")
}

func TestGetEmailByID_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)
	emailID := 1

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(emailID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sender_email", "recipient_email", "title", "isread", "sending_date", "description"}).
			AddRow("not_an_int", nil, nil, nil, nil, nil, nil))

	_, err = repo.GetEmailByID(emailID)
	assert.Error(t, err)
}

func TestInbox_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)
	userEmail := "test@example.com"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(userEmail).
		WillReturnRows(sqlmock.NewRows([]string{"sender_email", "sending_date", "isread", "title", "description"}).
			AddRow(nil, "not_a_time", nil, nil, nil))

	_, err = repo.Inbox(userEmail)
	assert.Error(t, err)
}

func TestGetEmailByID_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)
	emailID := 1

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(emailID).
		WillReturnError(errors.New("ошибка запроса к БД"))

	_, err = repo.GetEmailByID(emailID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка запроса к БД")
}

func TestInbox_EmptyEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)
	userEmail := ""

	query := `SELECT email_transaction.sender_email, email_transaction.sending_date, email_transaction.isread, message.title, message.description 
	FROM email_transaction AS t
	JOIN message AS m ON t.message_id = m.id
	WHERE t.recipient_email = $1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userEmail).
		WillReturnError(errors.New("email cannot be empty"))

	_, err = repo.Inbox(userEmail)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email cannot be empty")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Остались невыполненные ожидания: %v", err)
	}
}

func TestGetEmailByID_InvalidIDZero(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewEmailRepositoryService(db)

	query := `SELECT t.id, t.sender_email, t.recipient_email, m.title, t.isread, t.sending_date, m.description 
	FROM email_transaction AS t 
	JOIN message AS m ON t.message_id = m.id 
	WHERE t.id = $1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(0).
		WillReturnError(errors.New("invalid email ID"))

	_, err = repo.GetEmailByID(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email ID")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Остались невыполненные ожидания: %v", err)
	}
}
