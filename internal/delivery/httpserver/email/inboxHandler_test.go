package email

import (
	"context"
	"encoding/json"
	"errors"
	"mail/internal/delivery/httpserver/email/mocks"
	"mail/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInboxHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	emails := []models.Email{
		{
			ID:           1,
			Sender_email: "sender@example.com",
			Title:        "Test Email 1",
			Description:  "Test Description 1",
		},
		{
			ID:           2,
			Sender_email: "sender@example.com",
			Title:        "Test Email 2",
			Description:  "Test Description 2",
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/inbox", nil)
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		Inbox("test@example.com").
		Return(emails, nil)

	router.InboxHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response []models.Email
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, len(emails), len(response))

	for i := range emails {
		assert.Equal(t, emails[i].ID, response[i].ID)
		assert.Equal(t, emails[i].Sender_email, response[i].Sender_email)
		assert.Equal(t, emails[i].Title, response[i].Title)
		assert.Equal(t, emails[i].Description, response[i].Description)
	}
}

func TestInboxHandler_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	req := httptest.NewRequest(http.MethodGet, "/inbox", nil)
	rr := httptest.NewRecorder()

	router.InboxHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "unauthorized", response.Body)
}

func TestInboxHandler_InboxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	req := httptest.NewRequest(http.MethodGet, "/inbox", nil)
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		Inbox("test@example.com").
		Return(nil, errors.New("database error"))

	router.InboxHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "database error", response.Body)
}
