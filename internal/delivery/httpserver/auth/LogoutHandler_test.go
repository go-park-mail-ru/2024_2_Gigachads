package auth

import (
	"context"
	"encoding/json"
	"mail/internal/delivery/httpserver/email/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthRouter_LogoutHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	router := NewAuthRouter(mockUserUseCase)

	t.Run("успешный выход", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/logout", nil)
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockUserUseCase.EXPECT().
			Logout(gomock.Any(), "test@example.com").
			Return(nil)

		router.LogoutHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		cookies := w.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session" {
				sessionCookie = cookie
			}
		}

		assert.NotNil(t, sessionCookie)
		assert.Equal(t, "", sessionCookie.Value)
		assert.Equal(t, -1, sessionCookie.MaxAge)
	})

	t.Run("неавторизованный запрос", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/logout", nil)
		w := httptest.NewRecorder()

		router.LogoutHandler(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusUnauthorized), response["status"])
		assert.Equal(t, "unauthorized", response["body"])
	})

	t.Run("ошибка при выходе", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/logout", nil)
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockUserUseCase.EXPECT().
			Logout(gomock.Any(), "test@example.com").
			Return(assert.AnError)

		router.LogoutHandler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusInternalServerError), response["status"])
		assert.Equal(t, assert.AnError.Error(), response["body"])
	})
}
