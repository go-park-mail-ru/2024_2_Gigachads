package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	proto "mail/gen/go/smtp"
	"mail/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MockEmailRepository struct {
	mock.Mock
}

func (m *MockEmailRepository) Inbox(email string) ([]models.Email, error) {
	args := m.Called(email)
	return args.Get(0).([]models.Email), args.Error(1)
}

func (m *MockEmailRepository) GetEmailByID(id int) (models.Email, error) {
	args := m.Called(id)
	return args.Get(0).(models.Email), args.Error(1)
}

func (m *MockEmailRepository) GetSentEmails(email string) ([]models.Email, error) {
	args := m.Called(email)
	return args.Get(0).([]models.Email), args.Error(1)
}

func (m *MockEmailRepository) SaveEmail(email models.Email) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockEmailRepository) ChangeStatus(id int, status bool) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockEmailRepository) DeleteEmails(email string, ids []int) error {
	args := m.Called(email, ids)
	return args.Error(0)
}

func (m *MockEmailRepository) GetFolders(email string) ([]string, error) {
	args := m.Called(email)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockEmailRepository) GetFolderEmails(email string, folder string) ([]models.Email, error) {
	args := m.Called(email, folder)
	return args.Get(0).([]models.Email), args.Error(1)
}

func (m *MockEmailRepository) CreateFolder(email string, folder string) error {
	args := m.Called(email, folder)
	return args.Error(0)
}

func (m *MockEmailRepository) DeleteFolder(email string, folder string) error {
	args := m.Called(email, folder)
	return args.Error(0)
}

func (m *MockEmailRepository) RenameFolder(email string, folder string, newFolder string) error {
	args := m.Called(email, folder, newFolder)
	return args.Error(0)
}

func (m *MockEmailRepository) CheckFolder(email string, folder string) (bool, error) {
	args := m.Called(email, folder)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmailRepository) ChangeEmailFolder(id int, email string, folder string) error {
	args := m.Called(id, email, folder)
	return args.Error(0)
}

func (m *MockEmailRepository) CreateDraft(email models.Email) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockEmailRepository) UpdateDraft(draft models.Draft) error {
	args := m.Called(draft)
	return args.Error(0)
}

func (m *MockEmailRepository) GetMessageFolder(id int) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

type MockSmtpPop3ServiceClient struct {
	mock.Mock
}

func (m *MockSmtpPop3ServiceClient) SendEmail(ctx context.Context, req *proto.SendEmailRequest, opts ...grpc.CallOption) (*proto.SendEmailReply, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*proto.SendEmailReply), args.Error(1)
}

func (m *MockSmtpPop3ServiceClient) ForwardEmail(ctx context.Context, req *proto.ForwardEmailRequest, opts ...grpc.CallOption) (*proto.ForwardEmailReply, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*proto.ForwardEmailReply), args.Error(1)
}

func (m *MockSmtpPop3ServiceClient) ReplyEmail(ctx context.Context, req *proto.ReplyEmailRequest, opts ...grpc.CallOption) (*proto.ReplyEmailReply, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*proto.ReplyEmailReply), args.Error(1)
}

func (m *MockSmtpPop3ServiceClient) FetchEmailsViaPOP3(ctx context.Context, req *proto.FetchEmailsViaPOP3Request, opts ...grpc.CallOption) (*proto.FetchEmailsViaPOP3Reply, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*proto.FetchEmailsViaPOP3Reply), args.Error(1)
}

func TestEmailService_Inbox(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	expectedEmails := []models.Email{{ID: 1}, {ID: 2}}
	mockRepo.On("GetFolderEmails", "test@test.com", "Входящие").Return(expectedEmails, nil)

	emails, err := service.Inbox("test@test.com")
	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
	mockRepo.AssertExpectations(t)
}

func TestEmailService_GetEmailByID(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	expectedEmail := models.Email{ID: 1}
	mockRepo.On("GetEmailByID", 1).Return(expectedEmail, nil)

	email, err := service.GetEmailByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmail, email)
	mockRepo.AssertExpectations(t)
}

func TestEmailService_DeleteEmails(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	tests := []struct {
		name      string
		email     string
		ids       []int
		folder    string
		expectErr bool
	}{
		{
			name:   "Delete from trash",
			email:  "test@test.com",
			ids:    []int{1},
			folder: "Корзина",
		},
		{
			name:   "Move to trash",
			email:  "test@test.com",
			ids:    []int{1},
			folder: "Входящие",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetMessageFolder", tt.ids[0]).Return(tt.folder, nil).Once()
			if tt.folder == "Корзина" {
				mockRepo.On("DeleteEmails", tt.email, tt.ids).Return(nil).Once()
			} else {
				mockRepo.On("ChangeEmailFolder", tt.ids[0], tt.email, "Корзина").Return(nil).Once()
			}

			err := service.DeleteEmails(tt.email, tt.ids)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmailService_SendEmail(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	ctx := context.Background()
	from := "sender@test.com"
	to := []string{"recipient1@test.com", "recipient2@test.com"}
	subject := "Test Subject"
	body := "Test Body"

	tests := []struct {
		name    string
		mockErr error
	}{
		{
			name:    "Success send",
			mockErr: nil,
		},
		{
			name:    "Error sending email",
			mockErr: fmt.Errorf("smtp error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// В случае ошибки, мы ожидаем только один вызов, так как функция должна прервать выполнение
			expectedCalls := len(to)
			if tt.mockErr != nil {
				expectedCalls = 1
			}

			for i := 0; i < expectedCalls; i++ {
				req := &proto.SendEmailRequest{
					From:    from,
					To:      to[i],
					Subject: subject,
					Body:    body,
				}
				mockClient.On("SendEmail", ctx, req).Return(&proto.SendEmailReply{}, tt.mockErr).Once()
			}

			err := service.SendEmail(ctx, from, to, subject, body)
			if tt.mockErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestEmailService_ForwardEmail(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	ctx := context.Background()
	from := "sender@test.com"
	to := []string{"recipient1@test.com", "recipient2@test.com"}
	originalEmail := models.Email{
		ID:           1,
		Sender_email: "original@test.com",
		Title:        "Original Subject",
		Description:  "Original Body",
		Sending_date: time.Now(),
	}

	tests := []struct {
		name    string
		mockErr error
	}{
		{
			name:    "Success forward",
			mockErr: nil,
		},
		{
			name:    "Error forwarding email",
			mockErr: fmt.Errorf("smtp error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// В случае ошибки, мы ожидаем только один вызов
			expectedCalls := len(to)
			if tt.mockErr != nil {
				expectedCalls = 1
			}

			for i := 0; i < expectedCalls; i++ {
				req := &proto.ForwardEmailRequest{
					From:        from,
					To:          to[i],
					Sender:      originalEmail.Sender_email,
					Title:       originalEmail.Title,
					Description: originalEmail.Description,
					SendingDate: timestamppb.New(originalEmail.Sending_date),
				}
				mockClient.On("ForwardEmail", ctx, req).Return(&proto.ForwardEmailReply{}, tt.mockErr).Once()
			}

			err := service.ForwardEmail(ctx, from, to, originalEmail)
			if tt.mockErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestEmailService_ReplyEmail(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	ctx := context.Background()
	from := "sender@test.com"
	to := "recipient@test.com"
	replyText := "Reply text"
	originalEmail := models.Email{
		ID:           1,
		Sender_email: "original@test.com",
		Title:        "Original Subject",
		Description:  "Original Body",
		Sending_date: time.Now(),
	}

	tests := []struct {
		name    string
		mockErr error
	}{
		{
			name:    "Success reply",
			mockErr: nil,
		},
		{
			name:    "Error replying to email",
			mockErr: fmt.Errorf("smtp error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &proto.ReplyEmailRequest{
				From:        from,
				To:          to,
				Sender:      originalEmail.Sender_email,
				Title:       originalEmail.Title,
				Description: originalEmail.Description,
				SendingDate: timestamppb.New(originalEmail.Sending_date),
				ReplyText:   replyText,
			}
			mockClient.On("ReplyEmail", ctx, req).Return(&proto.ReplyEmailReply{}, tt.mockErr).Once()

			err := service.ReplyEmail(ctx, from, to, originalEmail, replyText)
			if tt.mockErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestEmailService_SendDraft(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	email := models.Email{
		ID:           1,
		Sender_email: "sender@test.com",
		Recipient:    "recipient@test.com",
		Title:        "Draft Subject",
		Description:  "Draft Body",
	}

	tests := []struct {
		name          string
		mockGetErr    error
		mockDeleteErr error
		mockSaveErr   error
		expectedError bool
	}{
		{
			name:          "Success send draft",
			mockGetErr:    nil,
			mockDeleteErr: nil,
			mockSaveErr:   nil,
			expectedError: false,
		},
		{
			name:          "Error getting draft",
			mockGetErr:    fmt.Errorf("db error"),
			mockDeleteErr: nil,
			mockSaveErr:   nil,
			expectedError: true,
		},
		{
			name:          "Error deleting draft",
			mockGetErr:    nil,
			mockDeleteErr: fmt.Errorf("db error"),
			mockSaveErr:   nil,
			expectedError: true,
		},
		{
			name:          "Error saving email",
			mockGetErr:    nil,
			mockDeleteErr: nil,
			mockSaveErr:   fmt.Errorf("db error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetEmailByID", email.ID).Return(email, tt.mockGetErr).Once()
			if tt.mockGetErr == nil {
				mockRepo.On("DeleteEmails", email.Sender_email, []int{email.ID}).Return(tt.mockDeleteErr).Once()
				if tt.mockDeleteErr == nil {
					mockRepo.On("SaveEmail", email).Return(tt.mockSaveErr).Once()
				}
			}

			err := service.SendDraft(email)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEmailService_Folders(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	email := "test@test.com"
	defaultFolders := []string{"Входящие", "Отправленные", "Спам", "Черновики", "Корзина"}

	t.Run("Existing folders", func(t *testing.T) {
		mockRepo.On("GetFolders", email).Return(defaultFolders, nil).Times(2)

		folders, err := service.GetFolders(email)
		assert.NoError(t, err)
		assert.Equal(t, defaultFolders, folders)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create default folders", func(t *testing.T) {
		mockRepo.On("GetFolders", email).Return([]string{}, nil).Once()

		for _, folder := range defaultFolders {
			mockRepo.On("CreateFolder", email, folder).Return(nil).Once()
		}

		mockRepo.On("GetFolders", email).Return(defaultFolders, nil).Once()

		folders, err := service.GetFolders(email)
		assert.NoError(t, err)
		assert.Equal(t, defaultFolders, folders)
		mockRepo.AssertExpectations(t)
	})
}

func TestEmailService_FolderOperations(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	t.Run("Create folder", func(t *testing.T) {
		mockRepo.On("CheckFolder", "test@test.com", "NewFolder").Return(false, nil)
		mockRepo.On("CreateFolder", "test@test.com", "NewFolder").Return(nil)

		err := service.CreateFolder("test@test.com", "NewFolder")
		assert.NoError(t, err)
	})

	t.Run("Delete folder", func(t *testing.T) {
		mockRepo.On("DeleteFolder", "test@test.com", "CustomFolder").Return(nil)

		err := service.DeleteFolder("test@test.com", "CustomFolder")
		assert.NoError(t, err)

		err = service.DeleteFolder("test@test.com", "Входящие")
		assert.Error(t, err)
	})

	t.Run("Rename folder", func(t *testing.T) {
		mockRepo.On("CheckFolder", "test@test.com", "NewName").Return(false, nil)
		mockRepo.On("RenameFolder", "test@test.com", "OldName", "NewName").Return(nil)

		err := service.RenameFolder("test@test.com", "OldName", "NewName")
		assert.NoError(t, err)

		err = service.RenameFolder("test@test.com", "Входящие", "NewName")
		assert.Error(t, err)
	})
}

func TestEmailService_GetSentEmails(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	expectedEmails := []models.Email{
		{ID: 1, Sender_email: "test@test.com", Recipient: "recipient1@test.com"},
		{ID: 2, Sender_email: "test@test.com", Recipient: "recipient2@test.com"},
	}

	mockRepo.On("GetFolderEmails", "test@test.com", "Отправленные").Return(expectedEmails, nil)

	emails, err := service.GetSentEmails("test@test.com")
	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
	mockRepo.AssertExpectations(t)
}

func TestEmailService_SaveEmail(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	email := models.Email{
		ID:           1,
		Sender_email: "test@test.com",
		Recipient:    "recipient@test.com",
		Title:        "Test Email",
		Description:  "Test Content",
		Sending_date: time.Now(),
	}

	mockRepo.On("SaveEmail", email).Return(nil)

	err := service.SaveEmail(email)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEmailService_ChangeStatus(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	tests := []struct {
		name      string
		id        int
		status    bool
		mockError error
	}{
		{
			name:      "Mark as read",
			id:        1,
			status:    true,
			mockError: nil,
		},
		{
			name:      "Mark as unread",
			id:        2,
			status:    false,
			mockError: nil,
		},
		{
			name:      "Error changing status",
			id:        3,
			status:    true,
			mockError: fmt.Errorf("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("ChangeStatus", tt.id, tt.status).Return(tt.mockError).Once()

			err := service.ChangeStatus(tt.id, tt.status)
			if tt.mockError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.mockError, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEmailService_GetFolderEmails(t *testing.T) {
	mockRepo := new(MockEmailRepository)
	mockClient := new(MockSmtpPop3ServiceClient)
	service := NewEmailService(mockRepo, mockClient)

	tests := []struct {
		name       string
		email      string
		folder     string
		mockEmails []models.Email
		expected   []models.Email
		mockError  error
	}{
		{
			name:   "Get inbox emails",
			email:  "test@test.com",
			folder: "Входящие",
			mockEmails: []models.Email{
				{ID: 1, Sender_email: "sender@test.com", Recipient: "test@test.com"},
				{ID: 2, Sender_email: "sender2@test.com", Recipient: "test@test.com"},
			},
			expected: []models.Email{
				{ID: 1, Sender_email: "sender@test.com", Recipient: "test@test.com"},
				{ID: 2, Sender_email: "sender2@test.com", Recipient: "test@test.com"},
			},
		},
		{
			name:   "Get sent emails",
			email:  "test@test.com",
			folder: "Отправленные",
			mockEmails: []models.Email{
				{ID: 1, Sender_email: "test@test.com", Recipient: "recipient@test.com"},
				{ID: 2, Sender_email: "test@test.com", Recipient: "recipient2@test.com"},
			},
			expected: []models.Email{
				{ID: 1, Sender_email: "recipient@test.com", Recipient: "test@test.com"},
				{ID: 2, Sender_email: "recipient2@test.com", Recipient: "test@test.com"},
			},
		},
		{
			name:   "Get draft emails",
			email:  "test@test.com",
			folder: "Черновики",
			mockEmails: []models.Email{
				{ID: 1, Sender_email: "test@test.com", Recipient: "draft@test.com"},
			},
			expected: []models.Email{
				{ID: 1, Sender_email: "draft@test.com", Recipient: "test@test.com"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetFolderEmails", tt.email, tt.folder).Return(tt.mockEmails, tt.mockError).Once()

			emails, err := service.GetFolderEmails(tt.email, tt.folder)
			if tt.mockError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.mockError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expected), len(emails))

				for i := range emails {
					assert.Equal(t, tt.expected[i].Sender_email, emails[i].Sender_email)
					assert.Equal(t, tt.expected[i].Recipient, emails[i].Recipient)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
