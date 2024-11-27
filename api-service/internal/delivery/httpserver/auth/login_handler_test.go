package auth

import (
	"bytes"
	"encoding/json"
	"mail/api-service/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
)

func TestAuthRouter_LoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUseCase := mocks.NewMockAuthUseCase(ctrl)
	router := NewAuthRouter(mockAuthUseCase)

	t.Run("успешный вход", func(t *testing.T) {
		user := models.User{
			Email:    "test@example.com",
			Password: "password123",
		}
		userJSON, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(userJSON))
		w := httptest.NewRecorder()

		expectedAvatar := "/avatars/test.jpg"
		expectedName := "Test User"
		expectedSession := "session123"
		expectedCSRF := "csrf123"

		mockAuthUseCase.EXPECT().
			Login(gomock.Any(), &user).
			Return(expectedAvatar, expectedName, expectedSession, expectedCSRF, nil)

		router.LoginHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response models.UserLogin
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, response.Email)
		assert.Equal(t, expectedName, response.Name)
		assert.Equal(t, expectedAvatar, response.AvatarURL)

		cookies := w.Result().Cookies()
		var sessionCookie, csrfCookie *http.Cookie
		for _, cookie := range cookies {
			switch cookie.Name {
			case "email":
				sessionCookie = cookie
			case "csrf":
				csrfCookie = cookie
			}
		}

		assert.NotNil(t, sessionCookie)
		assert.Equal(t, expectedSession, sessionCookie.Value)
		assert.NotNil(t, csrfCookie)
		assert.Equal(t, expectedCSRF, csrfCookie.Value)
	})

	t.Run("неверный JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()

		router.LoginHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_json", response["body"])
	})

	t.Run("неверный email", func(t *testing.T) {
		user := models.User{
			Email:    "invalid-email",
			Password: "password123",
		}
		userJSON, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(userJSON))
		w := httptest.NewRecorder()

		router.LoginHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_input", response["body"])
	})

	t.Run("ошибка авторизации", func(t *testing.T) {
		user := models.User{
			Email:    "test@example.com",
			Password: "password123",
		}
		userJSON, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(userJSON))
		w := httptest.NewRecorder()

		mockAuthUseCase.EXPECT().
			Login(gomock.Any(), &user).
			Return("", "", "", "", assert.AnError)

		router.LoginHandler(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_login_or_password", response["body"])
	})

	t.Run("пустой email", func(t *testing.T) {
		user := models.User{
			Email:    "",
			Password: "password123",
		}
		userJSON, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(userJSON))
		w := httptest.NewRecorder()

		router.LoginHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_input", response["body"])
	})
	t.Run("пустое тело запроса", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", nil)
		w := httptest.NewRecorder()

		router.LoginHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_json", response["body"])
	})

	t.Run("ошибка установки cookie", func(t *testing.T) {
		user := models.User{
			Email:    "test@example.com",
			Password: "password123",
		}
		userJSON, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(userJSON))
		w := httptest.NewRecorder()

		mockAuthUseCase.EXPECT().
			Login(gomock.Any(), &user).
			Return("/avatars/test.jpg", "Test User", "session123", "csrf123", nil)

		router.LoginHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		cookies := w.Result().Cookies()
		assert.Len(t, cookies, 2)
	})
}

func TestAuthRouter_LoginHandler_AdditionalCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUseCase := mocks.NewMockAuthUseCase(ctrl)
	router := NewAuthRouter(mockAuthUseCase)

	t.Run("ошибка сервиса авторизации с пустыми значениями", func(t *testing.T) {
		user := models.User{
			Email:    "test@example.com",
			Password: "password123",
		}
		userJSON, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(userJSON))
		w := httptest.NewRecorder()

		mockAuthUseCase.EXPECT().
			Login(gomock.Any(), &user).
			Return("", "", "", "", assert.AnError)

		router.LoginHandler(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_login_or_password", response["body"])
	})
}
