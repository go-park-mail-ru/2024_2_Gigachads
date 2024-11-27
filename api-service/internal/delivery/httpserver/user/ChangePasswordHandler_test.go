package user

import (
	"bytes"
	"context"
	"encoding/json"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserRouter_ChangePasswordHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	router := NewUserRouter(mockUserUseCase)

	t.Run("успешная смена пароля", func(t *testing.T) {
		changePass := &models.ChangePassword{
			Password:    "newPassword123",
			RePassword:  "newPassword123",
			OldPassword: "oldPassword123",
		}
		body, _ := json.Marshal(changePass)

		req := httptest.NewRequest("POST", "/change-password", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockUserUseCase.EXPECT().
			ChangePassword("test@example.com", "newPassword123").
			Return(nil)

		router.ChangePasswordHandler(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("неавторизованный запрос", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/change-password", nil)
		w := httptest.NewRecorder()

		router.ChangePasswordHandler(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("невалидный JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/change-password", bytes.NewBuffer([]byte("invalid json")))
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		router.ChangePasswordHandler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("пароли не совпадают", func(t *testing.T) {
		changePass := &models.ChangePassword{
			Password:    "newPassword123",
			RePassword:  "differentPassword123",
			OldPassword: "oldPassword123",
		}
		body, _ := json.Marshal(changePass)

		req := httptest.NewRequest("POST", "/change-password", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		router.ChangePasswordHandler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
