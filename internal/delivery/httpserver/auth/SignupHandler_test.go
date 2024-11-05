package auth

import (
	"bytes"
	"encoding/json"
	"mail/internal/delivery/httpserver/email/mocks"
	"mail/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthRouter_SignupHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	router := NewAuthRouter(mockUserUseCase)

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

		session := &models.Session{
			ID:   "session_id",
			Name: "session",
			Time: time.Now().Add(24 * time.Hour),
		}
		csrf := &models.Csrf{
			ID:   "csrf_id",
			Name: "csrf",
			Time: time.Now().Add(24 * time.Hour),
		}

		mockUserUseCase.EXPECT().
			Signup(gomock.Any(), gomock.Any()).
			Return(user, session, csrf, nil)

		router.SignupHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		cookies := w.Result().Cookies()
		var sessionCookie, csrfCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session" {
				sessionCookie = cookie
			}
			if cookie.Name == "csrf" {
				csrfCookie = cookie
			}
		}

		assert.NotNil(t, sessionCookie)
		assert.NotNil(t, csrfCookie)
		assert.Equal(t, "session_id", sessionCookie.Value)
		assert.Equal(t, "csrf_id", csrfCookie.Value)
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

		mockUserUseCase.EXPECT().
			Signup(gomock.Any(), gomock.Any()).
			Return(nil, nil, nil, assert.AnError)

		router.SignupHandler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusInternalServerError), response["status"])
		assert.Equal(t, assert.AnError.Error(), response["body"])
	})
}
