package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/api-service/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEmailService_UploadAttach(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	ctx := context.Background()
	fileContent := []byte("test file content")
	filename := "test.txt"
	expectedPath := "/path/to/file/test.txt"

	tests := []struct {
		name          string
		setupMock     func()
		expectedPath  string
		expectedError error
	}{
		{
			name: "Успешная загрузка файла",
			setupMock: func() {
				mockRepo.EXPECT().
					UploadAttach(ctx, fileContent, filename).
					Return(expectedPath, nil)
			},
			expectedPath:  expectedPath,
			expectedError: nil,
		},
		{
			name: "Ошибка при загрузке",
			setupMock: func() {
				mockRepo.EXPECT().
					UploadAttach(ctx, fileContent, filename).
					Return("", errors.New("ошибка загрузки"))
			},
			expectedPath:  "",
			expectedError: errors.New("ошибка загрузки"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			path, err := emailService.UploadAttach(ctx, fileContent, filename)

			assert.Equal(t, tt.expectedPath, path)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_GetAttach(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	ctx := context.Background()
	path := "/path/to/file/test.txt"
	expectedContent := []byte("test file content")

	tests := []struct {
		name            string
		setupMock       func()
		expectedContent []byte
		expectedError   error
	}{
		{
			name: "Успешное получение файла",
			setupMock: func() {
				mockRepo.EXPECT().
					GetAttach(ctx, path).
					Return(expectedContent, nil)
			},
			expectedContent: expectedContent,
			expectedError:   nil,
		},
		{
			name: "Ошибка при получениыи файла",
			setupMock: func() {
				mockRepo.EXPECT().
					GetAttach(ctx, path).
					Return(nil, errors.New("файл не найден"))
			},
			expectedContent: nil,
			expectedError:   errors.New("файл не найден"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			content, err := emailService.GetAttach(ctx, path)

			assert.Equal(t, tt.expectedContent, content)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_DeleteAttach(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	ctx := context.Background()
	path := "/path/to/file/test.txt"

	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
	}{
		{
			name: "Успешное удаление файла",
			setupMock: func() {
				mockRepo.EXPECT().
					DeleteAttach(ctx, path).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Ошибка при удалении файла",
			setupMock: func() {
				mockRepo.EXPECT().
					DeleteAttach(ctx, path).
					Return(errors.New("ошибка удаления"))
			},
			expectedError: errors.New("ошибка удаления"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := emailService.DeleteAttach(ctx, path)

			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_InboxStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	ctx := context.Background()
	email := "test@mail.ru"
	frontLastModified := time.Now().Add(-time.Hour)
	lastModified := time.Now()
	expectedEmails := []models.Email{{ID: 1}, {ID: 2}}

	tests := []struct {
		name           string
		setupMock      func()
		expectedEmails []models.Email
		expectedError  error
	}{
		{
			name: "Есть новые письма",
			setupMock: func() {
				mockRepo.EXPECT().GetTimestamp(ctx, email).Return(lastModified, nil)
				mockRepo.EXPECT().GetNewEmails(email, frontLastModified).Return(expectedEmails, nil)
			},
			expectedEmails: expectedEmails,
			expectedError:  nil,
		},
		{
			name: "Нет новых писем",
			setupMock: func() {
				mockRepo.EXPECT().GetTimestamp(ctx, email).Return(frontLastModified, nil)
			},
			expectedEmails: nil,
			expectedError:  fmt.Errorf("not_modified"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			emails, err := emailService.InboxStatus(ctx, email, frontLastModified)

			assert.Equal(t, tt.expectedEmails, emails)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_CreateFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	email := "test@mail.ru"
	folderName := "Новая папка"

	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
	}{
		{
			name: "Успешное создание папки",
			setupMock: func() {
				mockRepo.EXPECT().CheckFolder(email, folderName).Return(false, nil)
				mockRepo.EXPECT().CreateFolder(email, folderName).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Папка уже существует",
			setupMock: func() {
				mockRepo.EXPECT().CheckFolder(email, folderName).Return(true, nil)
			},
			expectedError: fmt.Errorf("folder_already_exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := emailService.CreateFolder(email, folderName)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_DeleteFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	email := "test@mail.ru"

	tests := []struct {
		name          string
		folderName    string
		setupMock     func()
		expectedError error
	}{
		{
			name:       "Успешное удаление папки",
			folderName: "Тестовая папка",
			setupMock: func() {
				mockRepo.EXPECT().DeleteFolder(email, "Тестовая папка").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:          "Попытка удаления системной папки",
			folderName:    "Входящие",
			expectedError: fmt.Errorf("unable_to_delete_folder"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}
			err := emailService.DeleteFolder(email, tt.folderName)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_RenameFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	email := "test@mail.ru"
	oldName := "Старое название"
	newName := "Новое название"

	tests := []struct {
		name          string
		oldName       string
		setupMock     func()
		expectedError error
	}{
		{
			name:    "Успешное переименование папки",
			oldName: oldName,
			setupMock: func() {
				mockRepo.EXPECT().CheckFolder(email, newName).Return(false, nil)
				mockRepo.EXPECT().RenameFolder(email, oldName, newName).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:          "Попытка переименования системной папки",
			oldName:       "Входящие",
			expectedError: fmt.Errorf("unable_to_rename_folder"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}
			err := emailService.RenameFolder(email, tt.oldName, newName)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_ChangeEmailFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	email := "test@mail.ru"
	messageID := 1
	folderName := "Новая папка"

	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
	}{
		{
			name: "Успешное перемещение письма",
			setupMock: func() {
				mockRepo.EXPECT().ChangeEmailFolder(messageID, email, folderName).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Ошибка при перемещении",
			setupMock: func() {
				mockRepo.EXPECT().ChangeEmailFolder(messageID, email, folderName).Return(fmt.Errorf("ошибка перемещения"))
			},
			expectedError: fmt.Errorf("ошибка перемещения"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := emailService.ChangeEmailFolder(messageID, email, folderName)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_CreateDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	draftEmail := models.Email{
		ID:           1,
		Sender_email: "test@mail.ru",
		Title:        "Черновик",
	}

	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
	}{
		{
			name: "Успешное создание черновика",
			setupMock: func() {
				mockRepo.EXPECT().CreateDraft(draftEmail).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := emailService.CreateDraft(draftEmail)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_UpdateDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	draftEmail := models.Email{
		ID:           1,
		Sender_email: "test@mail.ru",
		Title:        "Обновленный черновик",
	}

	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
	}{
		{
			name: "Успешное обновление черновика",
			setupMock: func() {
				mockRepo.EXPECT().UpdateDraft(draftEmail).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := emailService.UpdateDraft(draftEmail)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_SendDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	draftEmail := models.Email{
		ID:           1,
		Sender_email: "test@mail.ru",
		Title:        "Черновик для отправки",
	}

	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
	}{
		{
			name: "Успешная отправка черновика",
			setupMock: func() {
				mockRepo.EXPECT().GetEmailByID(draftEmail.ID).Return(draftEmail, nil)
				mockRepo.EXPECT().DeleteEmails(draftEmail.Sender_email, []int{draftEmail.ID}).Return(nil)
				mockRepo.EXPECT().SaveEmail(draftEmail).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Ошибка при получении черновика",
			setupMock: func() {
				mockRepo.EXPECT().GetEmailByID(draftEmail.ID).Return(models.Email{}, fmt.Errorf("черновик не найден"))
			},
			expectedError: fmt.Errorf("черновик не найден"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := emailService.SendDraft(draftEmail)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_ChangeStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	tests := []struct {
		name          string
		id            int
		status        bool
		setupMock     func()
		expectedError error
	}{
		{
			name:   "Успешное изменение статуса на прочитанное",
			id:     1,
			status: true,
			setupMock: func() {
				mockRepo.EXPECT().ChangeStatus(1, true).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "Ошибка при изменении статуса",
			id:     2,
			status: false,
			setupMock: func() {
				mockRepo.EXPECT().ChangeStatus(2, false).Return(fmt.Errorf("ошибка изменения статуса"))
			},
			expectedError: fmt.Errorf("ошибка изменения статуса"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := emailService.ChangeStatus(tt.id, tt.status)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_DeleteEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	userEmail := "test@mail.ru"
	messageIDs := []int{1, 2}

	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
	}{
		{
			name: "Успешное удаление писем из корзины",
			setupMock: func() {
				mockRepo.EXPECT().GetMessageFolder(messageIDs[0]).Return("Корзина", nil)
				mockRepo.EXPECT().DeleteEmails(userEmail, messageIDs).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Перемещение писем в корзину",
			setupMock: func() {
				mockRepo.EXPECT().GetMessageFolder(messageIDs[0]).Return("Входящие", nil)
				mockRepo.EXPECT().ChangeEmailFolder(messageIDs[0], userEmail, "Корзина").Return(nil)
				mockRepo.EXPECT().GetMessageFolder(messageIDs[1]).Return("Входящие", nil)
				mockRepo.EXPECT().ChangeEmailFolder(messageIDs[1], userEmail, "Корзина").Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Ошибка при получении папки",
			setupMock: func() {
				mockRepo.EXPECT().GetMessageFolder(messageIDs[0]).Return("", fmt.Errorf("ошибка получения папки"))
			},
			expectedError: fmt.Errorf("ошибка получения папки"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := emailService.DeleteEmails(userEmail, messageIDs)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_GetFolders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	email := "test@mail.ru"
	existingFolders := []string{"Входящие", "Отправленные", "Спам", "Черновики", "Корзина"}

	tests := []struct {
		name            string
		setupMock       func()
		expectedFolders []string
		expectedError   error
	}{
		{
			name: "Создание системных папок для нового пользователя",
			setupMock: func() {
				mockRepo.EXPECT().GetFolders(email).Return([]string{}, nil)
				mockRepo.EXPECT().CreateFolder(email, "Входящие").Return(nil)
				mockRepo.EXPECT().CreateFolder(email, "Отправленные").Return(nil)
				mockRepo.EXPECT().CreateFolder(email, "Спам").Return(nil)
				mockRepo.EXPECT().CreateFolder(email, "Черновики").Return(nil)
				mockRepo.EXPECT().CreateFolder(email, "Корзина").Return(nil)
				mockRepo.EXPECT().GetFolders(email).Return(existingFolders, nil)
			},
			expectedFolders: existingFolders,
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			folders, err := emailService.GetFolders(email)
			assert.Equal(t, tt.expectedFolders, folders)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_GetFolderEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	email := "test@mail.ru"
	emails := []models.Email{
		{
			ID:           1,
			Sender_email: "sender@mail.ru",
			Recipient:    "test@mail.ru",
			Title:        "Test Email",
		},
	}

	sentEmails := []models.Email{
		{
			ID:           2,
			Sender_email: "test@mail.ru",
			Recipient:    "recipient@mail.ru",
			Title:        "Sent Email",
		},
	}

	tests := []struct {
		name           string
		folderName     string
		setupMock      func()
		expectedEmails []models.Email
		expectedError  error
	}{
		{
			name:       "Получение писем из папки Входящие",
			folderName: "Входящие",
			setupMock: func() {
				mockRepo.EXPECT().GetFolderEmails(email, "Входящие").Return(emails, nil)
			},
			expectedEmails: emails,
			expectedError:  nil,
		},
		{
			name:       "Получение писем из папки Отправленные",
			folderName: "Отправленные",
			setupMock: func() {
				mockRepo.EXPECT().GetFolderEmails(email, "Отправленные").Return(sentEmails, nil)
			},
			expectedEmails: []models.Email{
				{
					ID:           2,
					Sender_email: "recipient@mail.ru",
					Recipient:    "test@mail.ru",
					Title:        "Sent Email",
				},
			},
			expectedError: nil,
		},
		{
			name:       "Ошибка при получении писем",
			folderName: "Входящие",
			setupMock: func() {
				mockRepo.EXPECT().GetFolderEmails(email, "Входящие").Return(nil, fmt.Errorf("ошибка получения писем"))
			},
			expectedEmails: nil,
			expectedError:  fmt.Errorf("ошибка получения писем"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			emails, err := emailService.GetFolderEmails(email, tt.folderName)
			assert.Equal(t, tt.expectedEmails, emails)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_Inbox(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	email := "test@mail.ru"
	expectedEmails := []models.Email{
		{
			ID:           1,
			Sender_email: "sender1@mail.ru",
			Recipient:    "test@mail.ru",
			Title:        "Test Email 1",
		},
		{
			ID:           2,
			Sender_email: "sender2@mail.ru",
			Recipient:    "test@mail.ru",
			Title:        "Test Email 2",
		},
	}

	tests := []struct {
		name           string
		setupMock      func()
		expectedEmails []models.Email
		expectedError  error
	}{
		{
			name: "Успешное получение входящих писем",
			setupMock: func() {
				mockRepo.EXPECT().GetFolderEmails(email, "Входящие").Return(expectedEmails, nil)
			},
			expectedEmails: expectedEmails,
			expectedError:  nil,
		},
		{
			name: "Ошибка при получении входящих писем",
			setupMock: func() {
				mockRepo.EXPECT().GetFolderEmails(email, "Входящие").Return(nil, fmt.Errorf("ошибка получения писем"))
			},
			expectedEmails: nil,
			expectedError:  fmt.Errorf("ошибка получения писем"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			emails, err := emailService.Inbox(email)
			assert.Equal(t, tt.expectedEmails, emails)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_GetEmailByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	emailID := 1
	expectedEmail := models.Email{
		ID:           1,
		Sender_email: "sender@mail.ru",
		Recipient:    "test@mail.ru",
		Title:        "Test Email",
		Attachments:  []string{"/path/to/file/test.txt"},
	}

	expectedEmailWithFiles := models.Email{
		ID:           1,
		Sender_email: "sender@mail.ru",
		Recipient:    "test@mail.ru",
		Title:        "Test Email",
		Attachments:  []string{"/path/to/file/test.txt"},
		Files: []models.File{
			{
				Path: "/path/to/file/test.txt",
				Name: "test.txt",
			},
		},
	}

	tests := []struct {
		name          string
		setupMock     func()
		expectedEmail models.Email
		expectedError error
	}{
		{
			name: "Успешное получение письма",
			setupMock: func() {
				mockRepo.EXPECT().GetEmailByID(emailID).Return(expectedEmail, nil)
			},
			expectedEmail: expectedEmailWithFiles,
			expectedError: nil,
		},
		{
			name: "Ошибка при получении письма",
			setupMock: func() {
				mockRepo.EXPECT().GetEmailByID(emailID).Return(models.Email{}, fmt.Errorf("письмо не найдено"))
			},
			expectedEmail: models.Email{},
			expectedError: fmt.Errorf("письмо не найдено"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			email, err := emailService.GetEmailByID(emailID)
			assert.Equal(t, tt.expectedEmail, email)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_GetSentEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	email := "test@mail.ru"
	expectedEmails := []models.Email{
		{
			ID:           1,
			Sender_email: "test@mail.ru",
			Recipient:    "recipient1@mail.ru",
			Title:        "Sent Email 1",
		},
		{
			ID:           2,
			Sender_email: "test@mail.ru",
			Recipient:    "recipient2@mail.ru",
			Title:        "Sent Email 2",
		},
	}

	tests := []struct {
		name           string
		setupMock      func()
		expectedEmails []models.Email
		expectedError  error
	}{
		{
			name: "Успешное получение отправленных писем",
			setupMock: func() {
				mockRepo.EXPECT().GetFolderEmails(email, "Отправленные").Return(expectedEmails, nil)
			},
			expectedEmails: expectedEmails,
			expectedError:  nil,
		},
		{
			name: "Ошибка при получении отправленных писем",
			setupMock: func() {
				mockRepo.EXPECT().GetFolderEmails(email, "Отправленные").Return(nil, fmt.Errorf("ошибка получения писем"))
			},
			expectedEmails: nil,
			expectedError:  fmt.Errorf("ошибка получения писем"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			emails, err := emailService.GetSentEmails(email)
			assert.Equal(t, tt.expectedEmails, emails)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestEmailService_SaveEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmailRepository(ctrl)
	emailService := EmailService{EmailRepo: mockRepo}

	ctx := context.Background()
	email := models.Email{
		ID:           1,
		Sender_email: "sender@mail.ru",
		Recipient:    "test@mail.ru",
		Title:        "Test Email",
	}

	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
	}{
		{
			name: "Успешное сохранение письма",
			setupMock: func() {
				mockRepo.EXPECT().SaveEmail(email).Return(nil)
				mockRepo.EXPECT().SetTimestamp(ctx, email.Recipient)
			},
			expectedError: nil,
		},
		{
			name: "Ошибка при сохранении письма",
			setupMock: func() {
				mockRepo.EXPECT().SaveEmail(email).Return(fmt.Errorf("ошибка сохранения"))
			},
			expectedError: fmt.Errorf("ошибка сохранения"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := emailService.SaveEmail(ctx, email)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
