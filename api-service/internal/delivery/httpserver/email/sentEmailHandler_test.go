package email

import (
	"context"
	"encoding/json"
	"errors"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSentEmailsHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	mockLogger := mocks.NewMockLogable(ctrl)
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

	req := httptest.NewRequest(http.MethodGet, "/sent", nil)
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		GetSentEmails("test@example.com").
		Return(emails, nil)

	mockLogger.EXPECT().
		Info(gomock.Any()).
		AnyTimes()

	router.SentEmailsHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response []models.Email
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, len(emails), len(response))
}

func TestSentEmailsHandler_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	mockLogger := mocks.NewMockLogable(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	req := httptest.NewRequest(http.MethodGet, "/sent", nil)
	rr := httptest.NewRecorder()

	mockLogger.EXPECT().
		Error(gomock.Any()).
		AnyTimes()

	router.SentEmailsHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "unauthorized", response.Body)
}

func TestSentEmailsHandler_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	mockLogger := mocks.NewMockLogable(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	req := httptest.NewRequest(http.MethodGet, "/sent", nil)
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		GetSentEmails("test@example.com").
		Return(nil, errors.New("database error"))

	mockLogger.EXPECT().
		Error(gomock.Any()).
		AnyTimes()

	router.SentEmailsHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "failed_to_get_sent_emails", response.Body)
}
