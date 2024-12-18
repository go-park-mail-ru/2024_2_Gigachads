package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/api-service/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// Вспомогательная функция для создания тестового Redis клиента
func newTestRedis(t *testing.T) (*redis.Client, func()) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Ошибка создания мока Redis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, func() {
		client.Close()
		mr.Close()
	}
}

func TestEmailRepositoryService_Inbox(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockLogable(ctrl)

	repo := NewEmailRepositoryService(db, &redis.Client{}, mockLogger)

	testTime := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "sender_email", "recipient_email",
		"title", "sending_date", "isread", "description",
	}).AddRow(
		1, "sender@test.com", "recipient@test.com",
		"Test Title", testTime, false, "Test Description",
	)

	tests := []struct {
		name    string
		email   string
		mock    func()
		want    []models.Email
		wantErr bool
	}{
		{
			name:  "успешное получение писем",
			email: "test@example.com",
			mock: func() {
				mock.ExpectQuery(`SELECT t.id, t.sender_email, t.recipient_email, m.title,`).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			want: []models.Email{
				{
					ID:           1,
					Sender_email: "sender@test.com",
					Recipient:    "recipient@test.com",
					Title:        "Test Title",
					Sending_date: testTime,
					IsRead:       false,
					Description:  "Test Description",
				},
			},
			wantErr: false,
		},
		{
			name:  "ошибка БД",
			email: "test@example.com",
			mock: func() {
				mock.ExpectQuery(`SELECT t.id, t.sender_email, t.recipient_email, m.title,`).
					WithArgs("test@example.com").
					WillReturnError(sql.ErrConnDone)
				mockLogger.EXPECT().Error(gomock.Any())
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.Inbox(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEmailRepositoryService_GetEmailByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockLogable(ctrl)

	repo := NewEmailRepositoryService(db, &redis.Client{}, mockLogger)

	testTime := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "parent_transaction_id", "sender_email", "recipient_email",
		"title", "isread", "sending_date", "description", "message_id",
	}).AddRow(
		1, "0", "sender@test.com", "recipient@test.com",
		"Test Title", false, testTime, "Test Description", 1,
	)

	attachRows := sqlmock.NewRows([]string{"url"}).
		AddRow("path/to/file1.pdf").
		AddRow("path/to/file2.jpg")

	tests := []struct {
		name    string
		id      int
		mock    func()
		want    models.Email
		wantErr bool
	}{
		{
			name: "успешное получение письма",
			id:   1,
			mock: func() {
				mock.ExpectQuery(`SELECT t.id, parent_transaction_id, t.sender_email`).
					WithArgs(1).
					WillReturnRows(rows)
				mock.ExpectQuery(`SELECT url`).
					WithArgs(1).
					WillReturnRows(attachRows)
			},
			want: models.Email{
				ID:           1,
				ParentID:     0,
				Sender_email: "sender@test.com",
				Recipient:    "recipient@test.com",
				Title:        "Test Title",
				IsRead:       false,
				Sending_date: testTime,
				Description:  "Test Description",
				Attachments:  []string{"path/to/file1.pdf", "path/to/file2.jpg"},
			},
			wantErr: false,
		},
		{
			name: "письмо не найдено",
			id:   999,
			mock: func() {
				mock.ExpectQuery(`SELECT t.id, parent_transaction_id, t.sender_email`).
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
				mockLogger.EXPECT().Error(gomock.Any())
			},
			want:    models.Email{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.GetEmailByID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEmailRepositoryService_GetFolders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockLogable(ctrl)

	repo := NewEmailRepositoryService(db, &redis.Client{}, mockLogger)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Inbox").
		AddRow("Sent").
		AddRow("Custom")

	tests := []struct {
		name    string
		email   string
		mock    func()
		want    []string
		wantErr bool
	}{
		{
			name:  "успешное получение папок",
			email: "test@example.com",
			mock: func() {
				mock.ExpectQuery(`SELECT f\.name FROM folder AS f JOIN profile AS p ON f\.user_id = p\.id WHERE p\.email = \$1 ORDER BY f\.id`).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			want:    []string{"Inbox", "Sent", "Custom"},
			wantErr: false,
		},
		{
			name:  "ошибка БД",
			email: "test@example.com",
			mock: func() {
				mock.ExpectQuery(`SELECT f\.name FROM folder AS f JOIN profile AS p ON f\.user_id = p\.id WHERE p\.email = \$1 ORDER BY f\.id`).
					WithArgs("test@example.com").
					WillReturnError(errors.New("database error"))
				mockLogger.EXPECT().Error(gomock.Any())
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.GetFolders(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmailRepositoryService_CreateFolder(t *testing.T) {
	opts := sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)
	db, mock, err := sqlmock.New(opts)
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockLogable(ctrl)

	repo := NewEmailRepositoryService(db, &redis.Client{}, mockLogger)

	tests := []struct {
		name       string
		email      string
		folderName string
		mock       func()
		wantErr    bool
	}{
		{
			name:       "ошибка при получении ID пользователя",
			email:      "test@example.com",
			folderName: "NewFolder",
			mock: func() {
				mock.ExpectQuery(`SELECT id FROM profile WHERE email = $1`).
					WithArgs("test@example.com").
					WillReturnError(sql.ErrNoRows)
				mockLogger.EXPECT().Error(gomock.Any())
			},
			wantErr: true,
		},
		{
			name:       "ошибка при создании папки",
			email:      "test@example.com",
			folderName: "NewFolder",
			mock: func() {
				userRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`SELECT id FROM profile WHERE email = $1`).
					WithArgs("test@example.com").
					WillReturnRows(userRows)

				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO folder (user_id, name) VALUES($1, $2)`).
					WithArgs(1, "NewFolder").
					WillReturnError(errors.New("database error"))
				mockLogger.EXPECT().Error(gomock.Any())
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:       "ошибка при начале транзакции",
			email:      "test@example.com",
			folderName: "NewFolder",
			mock: func() {
				userRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`SELECT id FROM profile WHERE email = $1`).
					WithArgs("test@example.com").
					WillReturnRows(userRows)

				mock.ExpectBegin().WillReturnError(errors.New("begin transaction error"))
				mockLogger.EXPECT().Error(gomock.Any())
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.CreateFolder(tt.email, tt.folderName)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmailRepositoryService_SaveEmail(t *testing.T) {
	opts := sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)
	db, mock, err := sqlmock.New(opts)
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockLogable(ctrl)

	repo := NewEmailRepositoryService(db, &redis.Client{}, mockLogger)

	testTime := time.Now()
	testEmail := models.Email{
		Sender_email: "sender@test.com",
		Recipient:    "recipient@test.com",
		Title:        "Test Title",
		Description:  "Test Description",
		Sending_date: testTime,
		ParentID:     0,
	}

	tests := []struct {
		name    string
		email   models.Email
		mock    func()
		wantErr bool
	}{
		{
			name:  "ошибка при получении ID отправителя",
			email: testEmail,
			mock: func() {
				mock.ExpectQuery(`SELECT id FROM profile WHERE email = $1`).
					WithArgs(testEmail.Sender_email).
					WillReturnError(sql.ErrNoRows)
				mockLogger.EXPECT().Error(gomock.Any())
			},
			wantErr: true,
		},
		{
			name:  "ошибка при получении ID получателя",
			email: testEmail,
			mock: func() {
				senderRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`SELECT id FROM profile WHERE email = $1`).
					WithArgs(testEmail.Sender_email).
					WillReturnRows(senderRows)

				mock.ExpectQuery(`SELECT id FROM profile WHERE email = $1`).
					WithArgs(testEmail.Recipient).
					WillReturnError(sql.ErrNoRows)
				mockLogger.EXPECT().Error(gomock.Any())
			},
			wantErr: true,
		},
		{
			name:  "ошибка при вставке сообщения",
			email: testEmail,
			mock: func() {
				// Успешные запросы ID
				senderRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`SELECT id FROM profile WHERE email = $1`).
					WithArgs(testEmail.Sender_email).
					WillReturnRows(senderRows)

				recipientRows := sqlmock.NewRows([]string{"id"}).AddRow(2)
				mock.ExpectQuery(`SELECT id FROM profile WHERE email = $1`).
					WithArgs(testEmail.Recipient).
					WillReturnRows(recipientRows)

				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO message (title, description) VALUES ($1, $2) RETURNING id`).
					WithArgs(testEmail.Title, testEmail.Description).
					WillReturnError(errors.New("database error"))
				mockLogger.EXPECT().Error(gomock.Any())
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.SaveEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmailRepositoryService_ChangeStatus(t *testing.T) {
	opts := sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)
	db, mock, err := sqlmock.New(opts)
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockLogable(ctrl)

	repo := NewEmailRepositoryService(db, &redis.Client{}, mockLogger)

	tests := []struct {
		name    string
		id      int
		status  bool
		mock    func()
		wantErr bool
	}{
		{
			name:   "письмо не найдено",
			id:     999,
			status: true,
			mock: func() {
				mock.ExpectExec(`UPDATE email_transaction SET isread = $1 WHERE id = $2`).
					WithArgs(true, 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mockLogger.EXPECT().Error(gomock.Any())
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.ChangeStatus(tt.id, tt.status)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmailRepositoryService_GetSentEmails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockLogable(ctrl)

	repo := NewEmailRepositoryService(db, &redis.Client{}, mockLogger)

	testTime := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "sender_email", "recipient_email",
		"title", "sending_date", "isread", "description",
	}).AddRow(
		1, "sender@test.com", "recipient@test.com",
		"Test Title", testTime, false, "Test Description",
	)

	tests := []struct {
		name    string
		email   string
		mock    func()
		want    []models.Email
		wantErr bool
	}{
		{
			name:  "успешное получение отправленных писем",
			email: "sender@test.com",
			mock: func() {
				mock.ExpectQuery(`SELECT t.id, t.sender_email, t.recipient_email`).
					WithArgs("sender@test.com").
					WillReturnRows(rows)
			},
			want: []models.Email{
				{
					ID:           1,
					Sender_email: "sender@test.com",
					Recipient:    "recipient@test.com",
					Title:        "Test Title",
					Sending_date: testTime,
					IsRead:       false,
					Description:  "Test Description",
				},
			},
			wantErr: false,
		},
		{
			name:  "ошибка БД",
			email: "sender@test.com",
			mock: func() {
				mock.ExpectQuery(`SELECT t.id, t.sender_email, t.recipient_email`).
					WithArgs("sender@test.com").
					WillReturnError(sql.ErrConnDone)
				mockLogger.EXPECT().Error(gomock.Any())
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.GetSentEmails(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEmailRepositoryService_DeleteEmails(t *testing.T) {
	opts := sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)
	db, mock, err := sqlmock.New(opts)
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockLogable(ctrl)

	repo := NewEmailRepositoryService(db, &redis.Client{}, mockLogger)

	tests := []struct {
		name    string
		email   string
		ids     []int
		mock    func()
		wantErr bool
	}{

		{
			name:  "ошибка при удалении",
			email: "test@example.com",
			ids:   []int{1, 2, 3},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM email_transaction WHERE id = ANY($1) AND (sender_email = $2 OR recipient_email = $2)`).
					WithArgs(pq.Array([]int{1, 2, 3}), "test@example.com").
					WillReturnError(errors.New("database error"))
				mockLogger.EXPECT().Error(gomock.Any())
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.DeleteEmails(tt.email, tt.ids)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmailRepositoryService_GetFolderEmails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockLogable(ctrl)

	repo := NewEmailRepositoryService(db, &redis.Client{}, mockLogger)

	testTime := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "sender_email", "recipient_email",
		"title", "sending_date", "isread", "description",
	}).AddRow(
		1, "sender@test.com", "recipient@test.com",
		"Test Title", testTime, false, "Test Description",
	)

	tests := []struct {
		name       string
		email      string
		folderName string
		mock       func()
		want       []models.Email
		wantErr    bool
	}{
		{
			name:       "успешное получение писем из папки",
			email:      "test@example.com",
			folderName: "Custom",
			mock: func() {
				mock.ExpectQuery(`SELECT t.id, t.sender_email, t.recipient_email`).
					WithArgs("test@example.com", "Custom").
					WillReturnRows(rows)
			},
			want: []models.Email{
				{
					ID:           1,
					Sender_email: "sender@test.com",
					Recipient:    "recipient@test.com",
					Title:        "Test Title",
					Sending_date: testTime,
					IsRead:       false,
					Description:  "Test Description",
				},
			},
			wantErr: false,
		},
		{
			name:       "ошибка БД",
			email:      "test@example.com",
			folderName: "Custom",
			mock: func() {
				mock.ExpectQuery(`SELECT t.id, t.sender_email, t.recipient_email`).
					WithArgs("test@example.com", "Custom").
					WillReturnError(sql.ErrConnDone)
				mockLogger.EXPECT().Error(gomock.Any())
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.GetFolderEmails(tt.email, tt.folderName)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEmailRepositoryService_DeleteAttach(t *testing.T) {
	opts := sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)
	db, mock, err := sqlmock.New(opts)
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockLogable(ctrl)

	repo := NewEmailRepositoryService(db, &redis.Client{}, mockLogger)
	ctx := context.Background()

	tests := []struct {
		name    string
		url     string
		mock    func()
		wantErr bool
	}{
		{
			name: "ошибка при удалении",
			url:  "test.pdf",
			mock: func() {
				mockLogger.EXPECT().Error(gomock.Any())
				mock.ExpectExec(`DELETE FROM attachment WHERE url = $1`).
					WithArgs("test.pdf").
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.DeleteAttach(ctx, tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmailRepositoryService_UploadAttach(t *testing.T) {
	opts := sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)
	db, mock, err := sqlmock.New(opts)
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockLogable(ctrl)

	repo := NewEmailRepositoryService(db, &redis.Client{}, mockLogger)
	ctx := context.Background()

	testData := []byte("test file content")

	tests := []struct {
		name    string
		data    []byte
		url     string
		mock    func()
		wantID  int
		wantErr bool
	}{
		{
			name: "ошибка при загрузке",
			data: testData,
			url:  "test.pdf",
			mock: func() {
				mockLogger.EXPECT().Error(gomock.Any())
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO attachment \(data, url\) VALUES \(\$1, \$2\) RETURNING id`).
					WithArgs(testData, "test.pdf").
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			gotID, err := repo.UploadAttach(ctx, tt.data, tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantID, gotID)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmailRepositoryService_ConnectAttachToMessage(t *testing.T) {
	opts := sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)
	db, mock, err := sqlmock.New(opts)
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockLogable(ctrl)

	repo := NewEmailRepositoryService(db, &redis.Client{}, mockLogger)

	tests := []struct {
		name      string
		messageID int
		url       string
		mock      func()
		wantErr   bool
	}{
		{
			name:      "ошибка при связывании",
			messageID: 1,
			url:       "test.pdf",
			mock: func() {
				mockLogger.EXPECT().Error(gomock.Any())
				mock.ExpectExec(`INSERT INTO message_attachment (message_id, url) VALUES ($1, $2)`).
					WithArgs(1, "test.pdf").
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.ConnectAttachToMessage(tt.messageID, tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmailRepositoryService_CreateDraft(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	redisClient, cleanup := newTestRedis(t)
	defer cleanup()

	mockLogger := mocks.NewMockLogable(gomock.NewController(t))
	repo := NewEmailRepositoryService(db, redisClient, mockLogger)

	email := models.Email{
		Sender_email: "test@mail.ru",
		Title:        "Draft",
	}

	tests := []struct {
		name        string
		email       models.Email
		mockSetup   func()
		expectedErr error
	}{
		{
			name:  "Успешное создание черновика",
			email: email,
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO message").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT id FROM profile").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT id FROM folder").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectExec("INSERT INTO email_transaction").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.CreateDraft(tt.email)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestEmailRepositoryService_UpdateDraft(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	redisClient, cleanup := newTestRedis(t)
	defer cleanup()

	mockLogger := mocks.NewMockLogable(gomock.NewController(t))
	repo := NewEmailRepositoryService(db, redisClient, mockLogger)

	email := models.Email{
		ID:           1,
		Title:        "Updated Draft",
		Description:  "Updated Description",
		Sender_email: "sender@mail.ru",
		Recipient:    "recipient@mail.ru",
	}

	tests := []struct {
		name        string
		email       models.Email
		mockSetup   func()
		expectedErr error
	}{
		{
			name:  "Успешное обновление черновика",
			email: email,
			mockSetup: func() {
				mock.ExpectBegin()

				// Получаем message_id
				mock.ExpectQuery("SELECT message_id FROM email_transaction WHERE id = $1").
					WithArgs(email.ID).
					WillReturnRows(sqlmock.NewRows([]string{"message_id"}).AddRow(1))

				// Обновляем сообщение
				mock.ExpectExec("UPDATE message SET title = $2, description = $3 WHERE id = $1").
					WithArgs(1, email.Title, email.Description).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Обновляем транзакцию
				mock.ExpectExec("UPDATE email_transaction SET recipient_email = $2 WHERE id = $1").
					WithArgs(email.ID, email.Recipient).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name:  "Ошибка при обновлении сообщения",
			email: email,
			mockSetup: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT message_id FROM email_transaction WHERE id = $1").
					WithArgs(email.ID).
					WillReturnRows(sqlmock.NewRows([]string{"message_id"}).AddRow(1))

				mock.ExpectExec("UPDATE message SET title = $2, description = $3 WHERE id = $1").
					WithArgs(1, email.Title, email.Description).
					WillReturnError(fmt.Errorf("ошибка обновления"))

				mockLogger.EXPECT().Error(gomock.Any()).Times(1)

				mock.ExpectRollback()
			},
			expectedErr: fmt.Errorf("ошибка обновления"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.UpdateDraft(tt.email)
			assert.Equal(t, tt.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
