package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/api-service/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthRouter_SignupHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUseCase := mocks.NewMockAuthUseCase(ctrl)
	router := NewAuthRouter(mockAuthUseCase)

	t.Run("невалидный JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer([]byte("invalid json")))
		w := httptest.NewRecorder()

		router.SignupHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusBadRequest), response["status"])
		assert.Equal(t, "invalid_json", response["body"])
	})

	t.Run("невалидный email", func(t *testing.T) {
		user := &models.User{
			Email:      "invalid-email",
			Name:       "TestUser",
			Password:   "password123",
			RePassword: "password123",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		router.SignupHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusBadRequest), response["status"])
		assert.Equal(t, "invalid_email", response["body"])
	})

	t.Run("пароли не совпадают", func(t *testing.T) {
		user := &models.User{
			Email:      "test@example.com",
			Name:       "TestUser",
			Password:   "password123",
			RePassword: "password456",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		router.SignupHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusBadRequest), response["status"])
		assert.Equal(t, "invalid_password", response["body"])
	})

	t.Run("успешная регистрация", func(t *testing.T) {
		user := &models.User{
			Email:      "test@example.com",
			Name:       "TestUser",
			Password:   "password123",
			RePassword: "password123",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		sessionID := "session_id"
		csrfID := "csrf_id"

		mockAuthUseCase.EXPECT().
			Signup(gomock.Any(), gomock.Any()).
			Return(sessionID, csrfID, nil)

		router.SignupHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		cookies := w.Result().Cookies()
		var sessionCookie, csrfCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "email" {
				sessionCookie = cookie
			}
			if cookie.Name == "csrf" {
				csrfCookie = cookie
			}
		}

		assert.NotNil(t, sessionCookie)
		assert.NotNil(t, csrfCookie)
		assert.Equal(t, sessionID, sessionCookie.Value)
		assert.Equal(t, csrfID, csrfCookie.Value)
	})

	t.Run("ошибка при создании пользователя", func(t *testing.T) {
		user := &models.User{
			Email:      "test@example.com",
			Name:       "TestUser",
			Password:   "password123",
			RePassword: "password123",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		mockAuthUseCase.EXPECT().
			Signup(gomock.Any(), gomock.Any()).
			Return("", "", errors.New("ошибка создания пользователя"))

		router.SignupHandler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusInternalServerError), response["status"])
		assert.Equal(t, "error_with_signup", response["body"])
	})
}
