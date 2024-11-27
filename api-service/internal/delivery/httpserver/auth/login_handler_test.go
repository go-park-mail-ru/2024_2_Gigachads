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

func TestAuthRouter_LoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	router := NewAuthRouter(mockUserUseCase)

	t.Run("успешный вход", func(t *testing.T) {
		loginData := &models.User{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(loginData)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		returnedUser := &models.User{
			Email:     "test@example.com",
			Name:      "TestUser",
			AvatarURL: "/custom/avatar.png",
		}
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
			Login(gomock.Any(), gomock.Any()).
			Return(returnedUser, session, csrf, nil)

		router.LoginHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response models.UserLogin
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", response.Email)
		assert.Equal(t, "TestUser", response.Name)
		assert.Equal(t, "/custom/avatar.png", response.AvatarURL)

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

	t.Run("успешный вход без аватара", func(t *testing.T) {
		loginData := &models.User{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(loginData)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		returnedUser := &models.User{
			Email:     "test@example.com",
			Name:      "TestUser",
			AvatarURL: "",
		}
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
			Login(gomock.Any(), gomock.Any()).
			Return(returnedUser, session, csrf, nil)

		router.LoginHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.UserLogin
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", response.Email)
		assert.Equal(t, "TestUser", response.Name)
		assert.Equal(t, "/icons/default.png", response.AvatarURL)
	})

	t.Run("невалидный JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid JSON")))
		w := httptest.NewRecorder()

		router.LoginHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
