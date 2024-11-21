package usecase

import (
	"errors"
	"mail/internal/delivery/httpserver/email/mocks"
	"mail/internal/models"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEmailService_SendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	mockSMTPRepo.EXPECT().
		SendEmail("from@example.com", []string{"to@example.com"}, "subject", "body").
		Return(nil)

	err := service.SendEmail("from@example.com", []string{"to@example.com"}, "subject", "body")
	assert.NoError(t, err)
}

func TestEmailService_ForwardEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	originalEmail := models.Email{
		Title:        "Original Subject",
		Sender_email: "original@example.com",
		Description:  "Original body",
		Sending_date: time.Now(),
	}

	mockSMTPRepo.EXPECT().
		SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	err := service.ForwardEmail("from@example.com", []string{"to@example.com"}, originalEmail)
	assert.NoError(t, err)
}

func TestEmailService_ReplyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	originalEmail := models.Email{
		Title:        "Original Subject",
		Sender_email: "original@example.com",
		Description:  "Original body",
		Sending_date: time.Now(),
	}

	mockSMTPRepo.EXPECT().
		SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	err := service.ReplyEmail("from@example.com", "to@example.com", originalEmail, "Reply text")
	assert.NoError(t, err)
}

func TestEmailService_GetEmailByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	email := models.Email{ID: 1, Title: "Test Email"}

	mockEmailRepo.EXPECT().
		GetEmailByID(1).
		Return(email, nil)

	result, err := service.GetEmailByID(1)
	assert.NoError(t, err)
	assert.Equal(t, email, result)
}

func TestEmailService_FetchEmailsViaPOP3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	gomock.InOrder(
		mockPOP3Repo.EXPECT().Connect().Return(nil),
		mockPOP3Repo.EXPECT().FetchEmails(mockEmailRepo).Return(nil),
		mockPOP3Repo.EXPECT().Quit().Return(nil),
	)

	err := service.FetchEmailsViaPOP3()
	assert.NoError(t, err)
}

func TestEmailService_GetSentEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	emails := []models.Email{{ID: 1, Title: "Test Email"}}

	mockEmailRepo.EXPECT().
		GetSentEmails("test@example.com").
		Return(emails, nil)

	result, err := service.GetSentEmails("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, emails, result)
}

func TestEmailService_SaveEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	email := models.Email{ID: 1, Title: "Test Email"}

	mockEmailRepo.EXPECT().
		SaveEmail(email).
		Return(nil)

	err := service.SaveEmail(email)
	assert.NoError(t, err)
}

func TestEmailService_ChangeStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	mockEmailRepo.EXPECT().
		ChangeStatus(1, true).
		Return(nil)

	err := service.ChangeStatus(1, true)
	assert.NoError(t, err)
}

func TestEmailService_Errors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	t.Run("SendEmail Error", func(t *testing.T) {
		mockSMTPRepo.EXPECT().
			SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(errors.New("smtp error"))

		err := service.SendEmail("from@example.com", []string{"to@example.com"}, "subject", "body")
		assert.Error(t, err)
	})

	t.Run("POP3 Connect Error", func(t *testing.T) {
		mockPOP3Repo.EXPECT().
			Connect().
			Return(errors.New("connection error"))

		err := service.FetchEmailsViaPOP3()
		assert.Error(t, err)
	})

	t.Run("GetEmailByID Error", func(t *testing.T) {
		mockEmailRepo.EXPECT().
			GetEmailByID(gomock.Any()).
			Return(models.Email{}, errors.New("not found"))

		_, err := service.GetEmailByID(1)
		assert.Error(t, err)
	})
}

func TestEmailService_ForwardEmail_WithAttachments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	originalEmail := models.Email{
		Title:        "Original Subject",
		Sender_email: "original@example.com",
		Description:  "Original body with attachments",
		Sending_date: time.Now(),
	}

	mockSMTPRepo.EXPECT().
		SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	err := service.ForwardEmail("from@example.com", []string{"to1@example.com", "to2@example.com"}, originalEmail)
	assert.NoError(t, err)
}

func TestEmailService_ReplyEmail_WithQuotedText(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	originalEmail := models.Email{
		Title:        "Original Subject",
		Sender_email: "original@example.com",
		Description:  "Original\nmultiline\nbody",
		Sending_date: time.Now(),
	}

	mockSMTPRepo.EXPECT().
		SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	err := service.ReplyEmail("from@example.com", "to@example.com", originalEmail, "Reply\nwith\nmultiple\nlines")
	assert.NoError(t, err)
}

func TestEmailService_FetchEmailsViaPOP3_WithErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	t.Run("FetchEmails Error", func(t *testing.T) {
		mockPOP3Repo.EXPECT().Connect().Return(nil)
		mockPOP3Repo.EXPECT().FetchEmails(mockEmailRepo).Return(errors.New("fetch error"))
		mockPOP3Repo.EXPECT().Quit().Return(nil)

		err := service.FetchEmailsViaPOP3()
		assert.Error(t, err)
		assert.Equal(t, "fetch error", err.Error())
	})

	t.Run("Connect Error", func(t *testing.T) {
		mockPOP3Repo.EXPECT().Connect().Return(errors.New("connect error"))

		err := service.FetchEmailsViaPOP3()
		assert.Error(t, err)
		assert.Equal(t, "connect error", err.Error())
	})

	t.Run("Success with Quit Error", func(t *testing.T) {
		mockPOP3Repo.EXPECT().Connect().Return(nil)
		mockPOP3Repo.EXPECT().FetchEmails(mockEmailRepo).Return(nil)
		mockPOP3Repo.EXPECT().Quit().Return(errors.New("quit error"))

		err := service.FetchEmailsViaPOP3()
		assert.NoError(t, err)
	})
}

func TestEmailService_GetSentEmails_WithFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	emails := []models.Email{
		{ID: 1, Title: "Test Email 1", Sender_email: "test@example.com"},
		{ID: 2, Title: "Test Email 2", Sender_email: "test@example.com"},
	}

	mockEmailRepo.EXPECT().
		GetSentEmails("test@example.com").
		Return(emails, nil)

	result, err := service.GetSentEmails("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, emails, result)
}

func TestEmailService_SaveEmail_WithValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	t.Run("Valid Email", func(t *testing.T) {
		email := models.Email{
			ID:           1,
			Title:        "Test Email",
			Sender_email: "sender@example.com",
			Description:  "Test body",
			Sending_date: time.Now(),
		}

		mockEmailRepo.EXPECT().
			SaveEmail(email).
			Return(nil)

		err := service.SaveEmail(email)
		assert.NoError(t, err)
	})

	t.Run("Save Error", func(t *testing.T) {
		email := models.Email{ID: 2}
		mockEmailRepo.EXPECT().
			SaveEmail(email).
			Return(errors.New("save error"))

		err := service.SaveEmail(email)
		assert.Error(t, err)
	})
}

func TestEmailService_ChangeStatus_WithValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	testCases := []struct {
		name     string
		id       int
		status   bool
		mockResp error
		wantErr  bool
	}{
		{
			name:     "Valid Status Change",
			id:       1,
			status:   true,
			mockResp: nil,
			wantErr:  false,
		},
		{
			name:     "Invalid ID",
			id:       -1,
			status:   true,
			mockResp: errors.New("invalid id"),
			wantErr:  true,
		},
		{
			name:     "Database Error",
			id:       1,
			status:   true,
			mockResp: errors.New("db error"),
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockEmailRepo.EXPECT().
				ChangeStatus(tc.id, tc.status).
				Return(tc.mockResp)

			err := service.ChangeStatus(tc.id, tc.status)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmailService_SendEmail_MultipleRecipients(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	recipients := []string{"to1@example.com", "to2@example.com", "to3@example.com"}
	mockSMTPRepo.EXPECT().
		SendEmail("from@example.com", recipients, "Test Subject", "Test Body").
		Return(nil)

	err := service.SendEmail("from@example.com", recipients, "Test Subject", "Test Body")
	assert.NoError(t, err)
}

func TestEmailService_SendEmail_EmptyRecipients(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	mockSMTPRepo.EXPECT().
		SendEmail("from@example.com", []string{}, "subject", "body").
		Return(errors.New("no recipients"))

	err := service.SendEmail("from@example.com", []string{}, "subject", "body")
	assert.Error(t, err)
}

func TestEmailService_ForwardEmail_ComplexEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	originalEmail := models.Email{
		Title:        "Meeting Notes",
		Sender_email: "boss@example.com",
		Description:  "Important meeting notes\nWith multiple lines\nAnd formatting",
		Sending_date: time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC),
	}

	mockSMTPRepo.EXPECT().
		SendEmail(
			"employee@example.com",
			[]string{"colleague1@example.com", "colleague2@example.com"},
			"Fwd: Meeting Notes",
			gomock.Any(),
		).Return(nil)

	err := service.ForwardEmail(
		"employee@example.com",
		[]string{"colleague1@example.com", "colleague2@example.com"},
		originalEmail,
	)
	assert.NoError(t, err)
}

func TestEmailService_ReplyEmail_WithFormattedText(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	originalEmail := models.Email{
		Title:        "Question about project",
		Sender_email: "client@example.com",
		Description:  "Hello,\n\nI have a question about...\n\nBest regards,\nClient",
		Sending_date: time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC),
	}

	replyText := "Hi Client,\n\nThank you for your question.\n\nBest regards,\nTeam"

	mockSMTPRepo.EXPECT().
		SendEmail(
			"support@example.com",
			[]string{"client@example.com"},
			"Re: Question about project",
			gomock.Any(),
		).Return(nil)

	err := service.ReplyEmail("support@example.com", "client@example.com", originalEmail, replyText)
	assert.NoError(t, err)
}

func TestEmailService_GetSentEmails_DetailedCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, nil, mockSMTPRepo, mockPOP3Repo)

	testCases := []struct {
		name        string
		senderEmail string
		emails      []models.Email
		wantErr     error
	}{
		{
			name:        "Multiple Emails",
			senderEmail: "sender@example.com",
			emails: []models.Email{
				{ID: 1, Title: "Email 1", Sender_email: "sender@example.com"},
				{ID: 2, Title: "Email 2", Sender_email: "sender@example.com"},
			},
			wantErr: nil,
		},
		{
			name:        "No Emails",
			senderEmail: "new@example.com",
			emails:      []models.Email{},
			wantErr:     nil,
		},
		{
			name:        "Error Case",
			senderEmail: "error@example.com",
			emails:      nil,
			wantErr:     errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockEmailRepo.EXPECT().
				GetSentEmails(tc.senderEmail).
				Return(tc.emails, tc.wantErr)

			emails, err := service.GetSentEmails(tc.senderEmail)
			if tc.wantErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.emails, emails)
			}
		})
	}
}

func TestEmailService_Inbox(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailRepo := mocks.NewMockEmailRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockSMTPRepo := mocks.NewMockSMTPRepository(ctrl)
	mockPOP3Repo := mocks.NewMockPOP3Repository(ctrl)

	service := NewEmailService(mockEmailRepo, mockSessionRepo, mockSMTPRepo, mockPOP3Repo)

	t.Run("успешное получение писем", func(t *testing.T) {
		emails := []models.Email{
			{
				ID:           1,
				Sender_email: "sender@example.com",
				Title:        "Test Email 1",
				Description:  "Test Description 1",
				Sending_date: time.Now(),
			},
			{
				ID:           2,
				Sender_email: "sender@example.com",
				Title:        "Test Email 2",
				Description:  "Test Description 2",
				Sending_date: time.Now(),
			},
		}

		mockEmailRepo.EXPECT().
			Inbox("user@example.com").
			Return(emails, nil)

		result, err := service.Inbox("user@example.com")
		assert.NoError(t, err)
		assert.Equal(t, emails, result)
	})

	t.Run("ошибка получения писем", func(t *testing.T) {
		mockEmailRepo.EXPECT().
			Inbox("user@example.com").
			Return(nil, errors.New("database error"))

		result, err := service.Inbox("user@example.com")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
