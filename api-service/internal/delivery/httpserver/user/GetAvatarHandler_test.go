package user

import (
	"context"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserRouter_GetAvatarHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	router := NewUserRouter(mockUserUseCase)

	t.Run("успешное получение PNG аватара", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/avatar", nil)
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockUserUseCase.EXPECT().
			GetAvatar("test@example.com").
			Return([]byte("fake-image-data"), "avatar.png", nil)

		router.GetAvatarHandler(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "image/png", w.Header().Get("Content-Type"))
	})

	t.Run("успешное получение JPEG аватара", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/avatar", nil)
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockUserUseCase.EXPECT().
			GetAvatar("test@example.com").
			Return([]byte("fake-image-data"), "avatar.jpg", nil)

		router.GetAvatarHandler(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "image/jpeg", w.Header().Get("Content-Type"))
	})

	t.Run("неавторизованный запрос", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/avatar", nil)
		w := httptest.NewRecorder()

		router.GetAvatarHandler(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
