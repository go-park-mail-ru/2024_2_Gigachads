package repository

import (
	"database/sql"
	"errors"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestEmailRepositoryService_DeleteEmails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()

	repo := NewEmailRepositoryService(db, mockLogger)

	tests := []struct {
		name       string
		email      string
		messageIDs []int
		mockFunc   func()
		wantErr    bool
	}{
		{
			name:       "Success delete",
			email:      "test@example.com",
			messageIDs: []int{1, 2, 3},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM email_transaction").
					WithArgs(pq.Array([]int{1, 2, 3})).
					WillReturnResult(sqlmock.NewResult(1, 3))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:       "Delete error",
			email:      "test@example.com",
			messageIDs: []int{1},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM email_transaction").
					WithArgs(pq.Array([]int{1})).
					WillReturnError(errors.New("delete error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := repo.DeleteEmails(tt.email, tt.messageIDs)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmailRepositoryService_GetFolderEmails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()

	repo := NewEmailRepositoryService(db, mockLogger)

	tests := []struct {
		name       string
		email      string
		folderName string
		mockFunc   func()
		want       []models.Email
		wantErr    bool
	}{
		{
			name:       "Success get folder emails",
			email:      "test@example.com",
			folderName: "Inbox",
			mockFunc: func() {
				testTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				rows := sqlmock.NewRows([]string{
					"id", "sender_email", "recipient_email",
					"title", "sending_date", "isread", "description",
				}).AddRow(
					1, "sender@example.com", "test@example.com",
					"Test Email", testTime, true, "Test Content",
				)
				mock.ExpectQuery("SELECT t.id, t.sender_email").
					WithArgs("test@example.com", "Inbox").
					WillReturnRows(rows)
			},
			want: []models.Email{
				{
					ID:           1,
					Sender_email: "sender@example.com",
					Recipient:    "test@example.com",
					Title:        "Test Email",
					Description:  "Test Content",
					IsRead:       true,
					Sending_date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := repo.GetFolderEmails(tt.email, tt.folderName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEmailRepositoryService_UpdateDraft(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()

	repo := NewEmailRepositoryService(db, mockLogger)

	tests := []struct {
		name     string
		draft    models.Draft
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Success update draft",
			draft: models.Draft{
				ID:          1,
				Title:       "Updated Draft",
				Description: "Updated Content",
				Recipient:   "recipient@example.com",
			},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT message_id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"message_id"}).AddRow(1))
				mock.ExpectExec("UPDATE message").
					WithArgs(1, "Updated Draft", "Updated Content").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE email_transaction").
					WithArgs(1, "recipient@example.com").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := repo.UpdateDraft(tt.draft)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmailRepositoryService_GetSentEmails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()

	repo := NewEmailRepositoryService(db, mockLogger)

	testTime := time.Date(2024, time.November, 27, 5, 38, 9, 0, time.Local)

	tests := []struct {
		name        string
		senderEmail string
		mockFunc    func()
		want        []models.Email
		wantErr     bool
	}{
		{
			name:        "Success get sent emails",
			senderEmail: "sender@example.com",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "sender_email", "recipient_email", "title", "sending_date", "isread", "description"}).
					AddRow(1, "sender@example.com", "recipient@example.com", "Test Email", testTime, false, "Test Content")
				mock.ExpectQuery("SELECT t.id, t.sender_email").
					WithArgs("sender@example.com").
					WillReturnRows(rows)
			},
			want: []models.Email{
				{
					ID:           1,
					ParentID:     0,
					Sender_email: "sender@example.com",
					Recipient:    "recipient@example.com",
					Title:        "Test Email",
					IsRead:       false,
					Sending_date: testTime,
					Description:  "Test Content",
				},
			},
			wantErr: false,
		},
		{
			name:        "Error getting sent emails",
			senderEmail: "sender@example.com",
			mockFunc: func() {
				mock.ExpectQuery("SELECT t.id, t.sender_email").
					WithArgs("sender@example.com").
					WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := repo.GetSentEmails(tt.senderEmail)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEmailRepositoryService_GetEmailByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()

	repo := NewEmailRepositoryService(db, mockLogger)

	testTime := time.Date(2024, time.November, 27, 5, 37, 32, 0, time.Local)

	tests := []struct {
		name     string
		id       int
		mockFunc func()
		want     models.Email
		wantErr  bool
	}{
		{
			name: "Success get email by ID",
			id:   1,
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "parent_transaction_id", "sender_email", "recipient_email", "title", "isread", "sending_date", "description"}).
					AddRow(1, nil, "sender@example.com", "recipient@example.com", "Test Email", false, testTime, "Test Content")
				mock.ExpectQuery("SELECT t.id, parent_transaction_id, t.sender_email, t.recipient_email, m.title, t.isread, t.sending_date, m.description").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: models.Email{
				ID:           1,
				ParentID:     0,
				Sender_email: "sender@example.com",
				Recipient:    "recipient@example.com",
				Title:        "Test Email",
				IsRead:       false,
				Sending_date: testTime,
				Description:  "Test Content",
			},
			wantErr: false,
		},
		{
			name: "Error getting email by ID",
			id:   1,
			mockFunc: func() {
				mock.ExpectQuery("SELECT t.id, parent_transaction_id, t.sender_email, t.recipient_email, m.title, t.isread, t.sending_date, m.description").
					WithArgs(1).
					WillReturnError(errors.New("db error"))
			},
			want:    models.Email{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := repo.GetEmailByID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEmailRepositoryService_SaveEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()

	repo := NewEmailRepositoryService(db, mockLogger)

	testTime := time.Date(2024, time.November, 27, 5, 38, 9, 0, time.Local)

	tests := []struct {
		name     string
		email    models.Email
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Error saving email",
			email: models.Email{
				Sender_email: "sender@example.com",
				Recipient:    "recipient@example.com",
				Title:        "Test Email",
				Description:  "Test Content",
				IsRead:       false,
				Sending_date: testTime,
			},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO message \(title, description\) VALUES \(\$1, \$2\) RETURNING id`).
					WithArgs("Test Email", "Test Content").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := repo.SaveEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmailRepositoryService_DeleteFolder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()

	repo := NewEmailRepositoryService(db, mockLogger)

	tests := []struct {
		name       string
		email      string
		folderName string
		mockFunc   func()
		wantErr    bool
	}{
		{
			name:       "Success delete folder",
			email:      "test@example.com",
			folderName: "TestFolder",
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id FROM profile").
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectExec("DELETE FROM folder").
					WithArgs(1, "TestFolder").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:       "Error user not found",
			email:      "test@example.com",
			folderName: "TestFolder",
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id FROM profile").
					WithArgs("test@example.com").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := repo.DeleteFolder(tt.email, tt.folderName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmailRepositoryService_RenameFolder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()

	repo := NewEmailRepositoryService(db, mockLogger)

	tests := []struct {
		name          string
		email         string
		folderName    string
		newFolderName string
		mockFunc      func()
		wantErr       bool
	}{
		{
			name:          "Success rename folder",
			email:         "test@example.com",
			folderName:    "OldName",
			newFolderName: "NewName",
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id FROM profile").
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT id FROM folder").
					WithArgs(1, "OldName").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
				mock.ExpectExec("UPDATE folder").
					WithArgs(2, "NewName").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := repo.RenameFolder(tt.email, tt.folderName, tt.newFolderName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmailRepositoryService_ChangeEmailFolder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()

	repo := NewEmailRepositoryService(db, mockLogger)

	tests := []struct {
		name       string
		id         int
		email      string
		folderName string
		mockFunc   func()
		wantErr    bool
	}{
		{
			name:       "Success change folder",
			id:         1,
			email:      "test@example.com",
			folderName: "NewFolder",
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id FROM profile").
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT id FROM folder").
					WithArgs(1, "NewFolder").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
				mock.ExpectExec("UPDATE email_transaction").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := repo.ChangeEmailFolder(tt.id, tt.email, tt.folderName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmailRepositoryService_CheckFolder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()

	repo := NewEmailRepositoryService(db, mockLogger)

	tests := []struct {
		name       string
		email      string
		folderName string
		mockFunc   func()
		want       bool
		wantErr    bool
	}{
		{
			name:       "Folder exists",
			email:      "test@example.com",
			folderName: "TestFolder",
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT f.id").
					WithArgs("test@example.com", "TestFolder").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			want:    true,
			wantErr: false,
		},
		{
			name:       "Folder does not exist",
			email:      "test@example.com",
			folderName: "NonExistentFolder",
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT f.id").
					WithArgs("test@example.com", "NonExistentFolder").
					WillReturnError(sql.ErrNoRows)
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := repo.CheckFolder(tt.email, tt.folderName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEmailRepositoryService_GetMessageFolder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()

	repo := NewEmailRepositoryService(db, mockLogger)

	tests := []struct {
		name     string
		msgID    int
		mockFunc func()
		want     string
		wantErr  bool
	}{
		{
			name:  "Success get folder",
			msgID: 1,
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT folder_id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"folder_id"}).AddRow(2))
				mock.ExpectQuery("SELECT name").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("TestFolder"))
				mock.ExpectCommit()
			},
			want:    "TestFolder",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := repo.GetMessageFolder(tt.msgID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEmailRepositoryService_CreateDraft(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()

	repo := NewEmailRepositoryService(db, mockLogger)

	testTime := time.Now()

	tests := []struct {
		name     string
		email    models.Email
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Success create draft",
			email: models.Email{
				Sender_email: "test@example.com",
				Recipient:    "recipient@example.com",
				Title:        "Draft Title",
				Description:  "Draft Content",
				Sending_date: testTime,
				IsRead:       false,
			},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO message").
					WithArgs("Draft Title", "Draft Content").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT id FROM profile").
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery(`SELECT id FROM folder WHERE user_id = \$1 AND name = 'Черновики'`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
				mock.ExpectExec("INSERT INTO email_transaction").
					WithArgs("test@example.com", "recipient@example.com", testTime, false, 1, nil, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Error creating message",
			email: models.Email{
				Sender_email: "test@example.com",
				Title:        "Draft Title",
				Description:  "Draft Content",
			},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO message").
					WithArgs("Draft Title", "Draft Content").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := repo.CreateDraft(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
