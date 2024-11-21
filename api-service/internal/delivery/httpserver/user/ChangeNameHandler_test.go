package user

import (
	"bytes"
	"context"
	"encoding/json"
	"mail/internal/delivery/httpserver/email/mocks"
	"mail/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserRouter_ChangeNameHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	router := NewUserRouter(mockUserUseCase)

	t.Run("успешное изменение имени", func(t *testing.T) {
		changeName := &models.ChangeName{
			Name: "NewUserName",
		}
		body, _ := json.Marshal(changeName)

		req := httptest.NewRequest("POST", "/change-name", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockUserUseCase.EXPECT().
			ChangeName(gomock.Eq("test@example.com"), gomock.Any()).
			Return(nil)

		router.ChangeNameHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})

	t.Run("неавторизованный запрос", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/change-name", nil)
		w := httptest.NewRecorder()

		router.ChangeNameHandler(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusUnauthorized), response["status"])
		assert.Equal(t, "unauthorized", response["body"])
	})

	t.Run("невалидный JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/change-name", bytes.NewBuffer([]byte("invalid json")))
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		router.ChangeNameHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusBadRequest), response["status"])
		assert.Equal(t, "invalid_json", response["body"])
	})
}
