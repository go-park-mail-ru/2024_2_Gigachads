package repository

import (
	"database/sql"
	"errors"
	"mail/models"
	"mail/pkg/logger"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Создаем мок логгера
type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...any)        {}
func (m *mockLogger) Error(msg string, args ...any)       {}
func (m *mockLogger) Debug(msg string, args ...any)       {}
func (m *mockLogger) Warn(msg string, args ...any)        {}
func (m *mockLogger) InitLogger() (logger.Logable, error) { return m, nil }

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *EmailRepositoryService) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать sqlmock: %v", err)
	}

	mockLogger := &mockLogger{}
	repo := NewEmailRepositoryService(db, mockLogger)
	return db, mock, repo
}

func TestSaveEmail(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	email := models.Email{
		Sender_email: "sender@example.com",
		Recipient:    "recipient@example.com",
		Title:        "Тестовое письмо",
		Description:  "Описание письма",
		Sending_date: time.Now(),
		IsRead:       false,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO message (title, description) VALUES ($1, $2) RETURNING id`)).
		WithArgs(email.Title, email.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO email_transaction (sender_email, recipient_email, sending_date, isread, message_id, parent_transaction_id) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(email.Sender_email, email.Recipient, email.Sending_date, email.IsRead, 1, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.SaveEmail(email)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSentEmails(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	senderEmail := "sender@example.com"
	query := `SELECT t.id, t.sender_email, t.recipient_email, m.title, 
			  t.sending_date, t.isread, m.description
			  FROM email_transaction AS t
			  JOIN message AS m ON t.message_id = m.id
			  WHERE t.sender_email = $1
			  ORDER BY t.sending_date DESC`

	rows := sqlmock.NewRows([]string{"id", "sender_email", "recipient_email", "title", "sending_date", "isread", "description"}).
		AddRow(1, senderEmail, "recipient@example.com", "Тестовое письмо", time.Now(), false, "Описание письма")

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(senderEmail).
		WillReturnRows(rows)

	emails, err := repo.GetSentEmails(senderEmail)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(emails))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSentEmails_Error(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	senderEmail := "sender@example.com"
	query := `SELECT t.id, t.sender_email, t.recipient_email, m.title, 
			  t.sending_date, t.isread, m.description
			  FROM email_transaction AS t
			  JOIN message AS m ON t.message_id = m.id
			  WHERE t.sender_email = $1
			  ORDER BY t.sending_date DESC`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(senderEmail).
		WillReturnError(errors.New("ошибка получения отправленных писем"))

	_, err := repo.GetSentEmails(senderEmail)
	assert.Error(t, err)
	assert.Equal(t, "ошибка получения отправленных писем", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveEmail_InsertMessageError(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

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

	err := repo.SaveEmail(email)
	if err == nil {
		t.Errorf("Ожидалась ошибка, но её не было")
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveEmail_ExecError(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

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
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO email_transaction (sender_email, recipient_email, sending_date, isread, message_id, parent_transaction_id) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(email.Sender_email, email.Recipient, email.Sending_date, email.IsRead, 1, nil).
		WillReturnError(errors.New("ошибка выполнения Exec"))
	mock.ExpectRollback()

	err := repo.SaveEmail(email)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка выполнения Exec")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetEmailByID(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	expectedEmail := models.Email{
		ID:           1,
		Sender_email: "sender@example.com",
		Recipient:    "recipient@example.com",
		Title:        "Тестовое письмо",
		Description:  "Описание письма",
		Sending_date: time.Now(),
		IsRead:       false,
	}

	query := `SELECT t.id, parent_transaction_id, t.sender_email, t.recipient_email, 
			  m.title, t.isread, t.sending_date, m.description 
			  FROM email_transaction AS t 
			  JOIN message AS m ON t.message_id = m.id 
			  WHERE m.id = $1`

	rows := sqlmock.NewRows([]string{
		"id", "parent_transaction_id", "sender_email", "recipient_email",
		"title", "isread", "sending_date", "description",
	}).AddRow(
		expectedEmail.ID, nil, expectedEmail.Sender_email,
		expectedEmail.Recipient, expectedEmail.Title,
		expectedEmail.IsRead, expectedEmail.Sending_date,
		expectedEmail.Description,
	)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnRows(rows)

	email, err := repo.GetEmailByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmail.ID, email.ID)
	assert.Equal(t, expectedEmail.Sender_email, email.Sender_email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetEmailByID_InvalidID(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	query := `SELECT t.id, parent_transaction_id, t.sender_email, t.recipient_email, 
			  m.title, t.isread, t.sending_date, m.description 
			  FROM email_transaction AS t 
			  JOIN message AS m ON t.message_id = m.id 
			  WHERE m.id = $1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(-1).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetEmailByID(-1)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows.Error(), err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetEmailByID_InvalidIDZero(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	query := `SELECT t.id, parent_transaction_id, t.sender_email, t.recipient_email, 
			  m.title, t.isread, t.sending_date, m.description 
			  FROM email_transaction AS t 
			  JOIN message AS m ON t.message_id = m.id 
			  WHERE m.id = $1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(0).
		WillReturnError(errors.New("invalid email ID"))

	_, err := repo.GetEmailByID(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email ID")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInbox(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	userEmail := "user@example.com"

	query := `SELECT t.id, t.sender_email, t.recipient_email, m.title, 
			  t.sending_date, t.isread, m.description
			  FROM email_transaction AS t
			  JOIN message AS m ON t.message_id = m.id
			  WHERE t.recipient_email = $1
			  ORDER BY t.sending_date DESC`

	rows := sqlmock.NewRows([]string{"id", "sender_email", "recipient_email", "title", "sending_date", "isread", "description"}).
		AddRow(1, "sender1@example.com", userEmail, "Письмо 1", time.Now(), false, "Описание письма 1").
		AddRow(2, "sender2@example.com", userEmail, "Письмо 2", time.Now(), true, "Описание письма 2")

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userEmail).
		WillReturnRows(rows)

	emails, err := repo.Inbox(userEmail)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(emails))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInbox_NoEmails(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	userEmail := "empty@example.com"

	query := `SELECT t.id, t.sender_email, t.recipient_email, m.title, 
			  t.sending_date, t.isread, m.description
			  FROM email_transaction AS t
			  JOIN message AS m ON t.message_id = m.id
			  WHERE t.recipient_email = $1
			  ORDER BY t.sending_date DESC`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userEmail).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sender_email", "recipient_email", "title", "sending_date", "isread", "description"}))

	emails, err := repo.Inbox(userEmail)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(emails))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInbox_QueryError(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	userEmail := "user@example.com"

	query := `SELECT t.id, t.sender_email, t.recipient_email, m.title, 
			  t.sending_date, t.isread, m.description
			  FROM email_transaction AS t
			  JOIN message AS m ON t.message_id = m.id
			  WHERE t.recipient_email = $1
			  ORDER BY t.sending_date DESC`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userEmail).
		WillReturnError(errors.New("ошибка запроса"))

	_, err := repo.Inbox(userEmail)
	assert.Error(t, err)
	assert.Equal(t, "ошибка запроса", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveEmail_BeginTxError(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	mock.ExpectBegin().WillReturnError(errors.New("ошибка начала транзакции"))

	email := models.Email{
		Sender_email: "sender@example.com",
		Recipient:    "recipient@example.com",
		Title:        "Test",
		Description:  "Test body",
	}

	err := repo.SaveEmail(email)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка начала транзакции")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveEmail_CommitError(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

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

	err := repo.SaveEmail(email)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка коммита")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetEmailByID_ScanError(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	emailID := 1

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(emailID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sender_email", "recipient_email", "title", "isread", "sending_date", "description"}).
			AddRow("not_an_int", nil, nil, nil, nil, nil, nil))

	_, err := repo.GetEmailByID(emailID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInbox_ScanError(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	userEmail := "test@example.com"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(userEmail).
		WillReturnRows(sqlmock.NewRows([]string{"sender_email", "sending_date", "isread", "title", "description"}).
			AddRow(nil, "not_a_time", nil, nil, nil))

	_, err := repo.Inbox(userEmail)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetEmailByID_QueryError(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	emailID := 1

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(emailID).
		WillReturnError(errors.New("ошибка запроса к БД"))

	_, err := repo.GetEmailByID(emailID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка запроса к БД")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInbox_EmptyEmail(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	userEmail := ""

	query := `SELECT t.id, t.sender_email, t.recipient_email, m.title, 
			  t.sending_date, t.isread, m.description
			  FROM email_transaction AS t
			  JOIN message AS m ON t.message_id = m.id
			  WHERE t.recipient_email = $1
			  ORDER BY t.sending_date DESC`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userEmail).
		WillReturnError(errors.New("email cannot be empty"))

	_, err := repo.Inbox(userEmail)
	assert.Error(t, err)
	assert.Equal(t, "email cannot be empty", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEmailRepository_ChangeStatus(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	testCases := []struct {
		name       string
		id         int
		status     bool
		mockSetup  func()
		wantError  bool
		errorMatch string
	}{
		{
			name:   "Success - Mark as true",
			id:     1,
			status: true,
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE email_transaction").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantError: false,
		},
		{
			name:   "Success - Mark as false",
			id:     2,
			status: false,
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE email_transaction").
					WithArgs(2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantError: false,
		},
		{
			name:   "Error - Database Error",
			id:     3,
			status: true,
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE email_transaction").
					WithArgs(3).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantError:  true,
			errorMatch: "database error",
		},
		{
			name:   "Error - No Rows Affected",
			id:     4,
			status: true,
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE email_transaction").
					WithArgs(4).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			wantError: false,
		},
		{
			name:   "Error - Invalid Status",
			id:     5,
			status: false,
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE email_transaction").
					WithArgs(5).
					WillReturnError(errors.New("invalid status"))
				mock.ExpectRollback()
			},
			wantError:  true,
			errorMatch: "invalid status",
		},
		{
			name:   "Error - Begin Transaction",
			id:     6,
			status: true,
			mockSetup: func() {
				mock.ExpectBegin().WillReturnError(errors.New("begin transaction error"))
			},
			wantError:  true,
			errorMatch: "begin transaction error",
		},
		{
			name:   "Error - Commit Transaction",
			id:     7,
			status: true,
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE email_transaction").
					WithArgs(7).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			wantError:  true,
			errorMatch: "commit error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			err := repo.ChangeStatus(tc.id, tc.status)

			if tc.wantError {
				assert.Error(t, err)
				if tc.errorMatch != "" {
					assert.Contains(t, err.Error(), tc.errorMatch)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEmailRepository_ChangeStatus_TransactionError(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	mock.ExpectBegin().WillReturnError(errors.New("transaction error"))

	err := repo.ChangeStatus(1, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
