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
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInboxHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fixedTime := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	expectedEmails := []models.Email{
		{
			ID:           1,
			Sender_email: "sender@example.com",
			Recipient:    "user@example.com",
			Title:        "Test Email 1",
			Description:  "Test Body 1",
			Sending_date: fixedTime,
			IsRead:       false,
		},
		{
			ID:           2,
			Sender_email: "another@example.com",
			Recipient:    "user@example.com",
			Title:        "Test Email 2",
			Description:  "Test Body 2",
			Sending_date: fixedTime,
			IsRead:       false,
		},
	}

	mockUseCase.EXPECT().
		Inbox("user@example.com").
		Return(expectedEmails, nil)

	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}

	req, err := http.NewRequest("GET", "/inbox", nil)
	assert.NoError(t, err)
	req = req.WithContext(context.WithValue(req.Context(), "email", "user@example.com"))
	rr := httptest.NewRecorder()

	emailRouter.InboxHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response []models.Email
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, response)
}

func TestInboxHandler_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}

	req, err := http.NewRequest("GET", "/inbox", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()

	emailRouter.InboxHandler(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestInboxHandler_InboxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	mockUseCase.EXPECT().
		Inbox("user@example.com").
		Return(nil, errors.New("database error"))

	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}

	req, err := http.NewRequest("GET", "/inbox", nil)
	assert.NoError(t, err)
	req = req.WithContext(context.WithValue(req.Context(), "email", "user@example.com"))
	rr := httptest.NewRecorder()

	emailRouter.InboxHandler(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestInboxHandler_EmptyInbox(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	mockUseCase.EXPECT().
		Inbox("user@example.com").
		Return([]models.Email{}, nil)

	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}

	req, err := http.NewRequest("GET", "/inbox", nil)
	assert.NoError(t, err)
	req = req.WithContext(context.WithValue(req.Context(), "email", "user@example.com"))
	rr := httptest.NewRecorder()

	emailRouter.InboxHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response []models.Email
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Empty(t, response)
}
