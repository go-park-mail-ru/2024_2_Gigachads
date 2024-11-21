package usecase

import (
	"context"
	"errors"
	"mail/internal/delivery/httpserver/email/mocks"
	"mail/internal/models"
	"mime/multipart"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockCsrfRepo := mocks.NewMockCsrfRepository(ctrl)
	ctx := context.Background()

	userService := NewUserService(mockUserRepo, mockSessionRepo, mockCsrfRepo)

	t.Run("успешная регистрация", func(t *testing.T) {
		user := &models.User{
			Email:    "test@example.com",
			Password: "password",
			Name:     "Test User",
		}

		mockUserRepo.EXPECT().
			IsExist(user.Email).
			Return(false, nil)

		mockUserRepo.EXPECT().
			CreateUser(user).
			Return(user, nil)

		mockSessionRepo.EXPECT().
			CreateSession(ctx, user.Email).
			Return(&models.Session{ID: "session-id"}, nil)

		mockCsrfRepo.EXPECT().
			CreateCsrf(ctx, user.Email).
			Return(&models.Csrf{ID: "csrf-id"}, nil)

		resultUser, session, csrf, err := userService.Signup(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, user, resultUser)
		assert.Equal(t, "session-id", session.ID)
		assert.Equal(t, "csrf-id", csrf.ID)
	})

	t.Run("пользователь уже существует", func(t *testing.T) {
		user := &models.User{Email: "existing@example.com"}

		mockUserRepo.EXPECT().
			IsExist(user.Email).
			Return(true, nil)

		_, _, _, err := userService.Signup(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, "login_taken", err.Error())
	})
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockCsrfRepo := mocks.NewMockCsrfRepository(ctrl)
	ctx := context.Background()

	userService := NewUserService(mockUserRepo, mockSessionRepo, mockCsrfRepo)

	t.Run("успешный вход", func(t *testing.T) {
		user := &models.User{
			Email:    "test@example.com",
			Password: "password",
		}

		mockUserRepo.EXPECT().
			IsExist(user.Email).
			Return(true, nil)

		mockUserRepo.EXPECT().
			CheckUser(user).
			Return(user, nil)

		mockSessionRepo.EXPECT().
			CreateSession(ctx, user.Email).
			Return(&models.Session{ID: "session-id"}, nil)

		mockCsrfRepo.EXPECT().
			CreateCsrf(ctx, user.Email).
			Return(&models.Csrf{
				ID:        "csrf-id",
				Name:      "csrf-name",
				Time:      time.Now(),
				UserLogin: user.Email,
			}, nil)

		resultUser, session, csrf, err := userService.Login(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, user, resultUser)
		assert.Equal(t, "session-id", session.ID)
		assert.Equal(t, "csrf-id", csrf.ID)
	})

	t.Run("пользователь не существует", func(t *testing.T) {
		user := &models.User{Email: "nonexistent@example.com"}

		mockUserRepo.EXPECT().
			IsExist(user.Email).
			Return(false, nil)

		resultUser, session, csrf, err := userService.Login(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, "user_does_not_exist", err.Error())
		assert.Nil(t, resultUser)
		assert.Nil(t, session)
		assert.Nil(t, csrf)
	})

	t.Run("ошибка проверки пользователя", func(t *testing.T) {
		user := &models.User{Email: "test@example.com"}

		mockUserRepo.EXPECT().
			IsExist(user.Email).
			Return(true, nil)

		mockUserRepo.EXPECT().
			CheckUser(user).
			Return(nil, errors.New("check user error"))

		resultUser, session, csrf, err := userService.Login(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, "check user error", err.Error())
		assert.Nil(t, resultUser)
		assert.Nil(t, session)
		assert.Nil(t, csrf)
	})
}

func TestUserService_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockCsrfRepo := mocks.NewMockCsrfRepository(ctrl)
	ctx := context.Background()

	userService := NewUserService(mockUserRepo, mockSessionRepo, mockCsrfRepo)

	t.Run("успешный выход", func(t *testing.T) {
		sessionID := "test-session"

		mockSessionRepo.EXPECT().
			DeleteSession(ctx, sessionID).
			Return(nil)

		mockCsrfRepo.EXPECT().
			DeleteCsrf(ctx, sessionID).
			Return(nil)

		err := userService.Logout(ctx, sessionID)
		assert.NoError(t, err)
	})

	t.Run("ошибка удаления сессии", func(t *testing.T) {
		sessionID := "test-session"

		mockSessionRepo.EXPECT().
			DeleteSession(ctx, sessionID).
			Return(errors.New("session error"))

		err := userService.Logout(ctx, sessionID)
		assert.Error(t, err)
	})
}

func TestUserService_CheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockCsrfRepo := mocks.NewMockCsrfRepository(ctrl)
	ctx := context.Background()

	userService := NewUserService(mockUserRepo, mockSessionRepo, mockCsrfRepo)

	t.Run("успешная проверка", func(t *testing.T) {
		sessionID := "test-session"
		email := "test@example.com"

		mockSessionRepo.EXPECT().
			GetSession(ctx, sessionID).
			Return(email, nil)

		resultEmail, err := userService.CheckAuth(ctx, sessionID)
		assert.NoError(t, err)
		assert.Equal(t, email, resultEmail)
	})

	t.Run("ошибка проверки", func(t *testing.T) {
		sessionID := "test-session"

		mockSessionRepo.EXPECT().
			GetSession(ctx, sessionID).
			Return("", errors.New("session error"))

		resultEmail, err := userService.CheckAuth(ctx, sessionID)
		assert.Error(t, err)
		assert.Empty(t, resultEmail)
	})
}

func TestUserService_CheckCsrf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockCsrfRepo := mocks.NewMockCsrfRepository(ctrl)
	ctx := context.Background()

	userService := NewUserService(mockUserRepo, mockSessionRepo, mockCsrfRepo)

	t.Run("успешная проверка CSRF", func(t *testing.T) {
		sessionID := "test-session"
		csrfToken := "test-csrf"
		email := "test@example.com"

		mockCsrfRepo.EXPECT().
			GetCsrf(ctx, csrfToken).
			Return(email, nil)

		mockSessionRepo.EXPECT().
			GetSession(ctx, sessionID).
			Return(email, nil)

		err := userService.CheckCsrf(ctx, sessionID, csrfToken)
		assert.NoError(t, err)
	})

	t.Run("несовпадающие email", func(t *testing.T) {
		sessionID := "test-session"
		csrfToken := "test-csrf"

		mockCsrfRepo.EXPECT().
			GetCsrf(ctx, csrfToken).
			Return("test1@example.com", nil)

		mockSessionRepo.EXPECT().
			GetSession(ctx, sessionID).
			Return("test2@example.com", nil)

		err := userService.CheckCsrf(ctx, sessionID, csrfToken)
		assert.Error(t, err)
		assert.Equal(t, "invalid_csrf", err.Error())
	})
}

func TestUserService_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockCsrfRepo := mocks.NewMockCsrfRepository(ctrl)

	userService := NewUserService(mockUserRepo, mockSessionRepo, mockCsrfRepo)

	t.Run("успешное изменение пароля", func(t *testing.T) {
		email := "test@example.com"
		newPassword := "newpassword"
		user := &models.User{
			Email:    email,
			Password: "oldpassword",
		}

		mockUserRepo.EXPECT().
			GetUserByEmail(email).
			Return(user, nil)

		mockUserRepo.EXPECT().
			UpdateInfo(gomock.Any()).
			Return(nil)

		err := userService.ChangePassword(email, newPassword)
		assert.NoError(t, err)
	})

	t.Run("ошибка получения пользователя", func(t *testing.T) {
		email := "test@example.com"
		newPassword := "newpassword"

		mockUserRepo.EXPECT().
			GetUserByEmail(email).
			Return(nil, errors.New("user not found"))

		err := userService.ChangePassword(email, newPassword)
		assert.Error(t, err)
	})
}

func TestUserService_ChangeName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockCsrfRepo := mocks.NewMockCsrfRepository(ctrl)

	userService := NewUserService(mockUserRepo, mockSessionRepo, mockCsrfRepo)

	t.Run("успешное изменение имени", func(t *testing.T) {
		email := "test@example.com"
		newName := "New Name"
		user := &models.User{
			Email: email,
			Name:  "Old Name",
		}

		mockUserRepo.EXPECT().
			GetUserByEmail(email).
			Return(user, nil)

		mockUserRepo.EXPECT().
			UpdateInfo(gomock.Any()).
			Return(nil)

		err := userService.ChangeName(email, newName)
		assert.NoError(t, err)
	})

	t.Run("ошибка получения пользователя", func(t *testing.T) {
		email := "test@example.com"
		newName := "New Name"

		mockUserRepo.EXPECT().
			GetUserByEmail(email).
			Return(nil, errors.New("user not found"))

		err := userService.ChangeName(email, newName)
		assert.Error(t, err)
	})
}

func TestUserService_ChangeAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockCsrfRepo := mocks.NewMockCsrfRepository(ctrl)

	userService := NewUserService(mockUserRepo, mockSessionRepo, mockCsrfRepo)

	t.Run("успешное изменение аватара", func(t *testing.T) {
		tmpFile := createTempFile(t, "test-avatar", []byte("test content"))
		defer os.Remove(tmpFile.Name())

		header := multipart.FileHeader{
			Filename: "test.jpg",
			Size:     100,
		}

		email := "test@example.com"
		user := &models.User{
			ID:    1,
			Email: email,
		}

		mockUserRepo.EXPECT().
			GetUserByEmail(email).
			Return(user, nil)

		mockUserRepo.EXPECT().
			UpdateInfo(gomock.Any()).
			Return(nil)

		err := userService.ChangeAvatar(tmpFile, header, email)
		assert.NoError(t, err)
	})

	t.Run("слишком большой файл", func(t *testing.T) {
		tmpFile := createTempFile(t, "test-avatar", []byte("test content"))
		defer os.Remove(tmpFile.Name())

		header := multipart.FileHeader{
			Filename: "test.jpg",
			Size:     6 * 1024 * 1024,
		}

		err := userService.ChangeAvatar(tmpFile, header, "test@example.com")
		assert.Error(t, err)
		assert.Equal(t, "too_big_file", err.Error())
	})
}

func TestUserService_GetAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockCsrfRepo := mocks.NewMockCsrfRepository(ctrl)

	userService := NewUserService(mockUserRepo, mockSessionRepo, mockCsrfRepo)

	t.Run("успешное получение аватара", func(t *testing.T) {
		err := os.MkdirAll("./avatars", os.ModePerm)
		assert.NoError(t, err)
		defer os.RemoveAll("./avatars")

		avatarContent := []byte("test avatar content")
		err = os.WriteFile("./avatars/test.jpg", avatarContent, 0644)
		assert.NoError(t, err)

		email := "test@example.com"
		user := &models.User{
			Email:     email,
			AvatarURL: "test.jpg",
		}

		mockUserRepo.EXPECT().
			GetUserByEmail(email).
			Return(user, nil)

		data, filename, err := userService.GetAvatar(email)
		assert.NoError(t, err)
		assert.Equal(t, avatarContent, data)
		assert.Equal(t, "test.jpg", filename)
	})

	t.Run("пользователь не найден", func(t *testing.T) {
		email := "test@example.com"

		mockUserRepo.EXPECT().
			GetUserByEmail(email).
			Return(nil, errors.New("user not found"))

		data, filename, err := userService.GetAvatar(email)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Empty(t, filename)
	})
}

func createTempFile(t *testing.T, prefix string, content []byte) *os.File {
	tmpFile, err := os.CreateTemp("", prefix)
	assert.NoError(t, err)

	_, err = tmpFile.Write(content)
	assert.NoError(t, err)

	return tmpFile
}
