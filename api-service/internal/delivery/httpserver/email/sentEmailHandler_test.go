package email

import (
	"context"
	"encoding/json"
	"mail/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"mail/internal/delivery/httpserver/email/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSentEmailsHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := &EmailRouter{EmailUseCase: mockEmailUseCase}

	senderEmail := "test@example.com"
	emails := []models.Email{
		{Sender_email: senderEmail, Recipient: "recipient@example.com", Title: "Test Email", Description: "This is a test email."},
	}

	req := httptest.NewRequest(http.MethodGet, "/sent-emails", nil)
	ctx := context.WithValue(req.Context(), "email", senderEmail)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().GetSentEmails(senderEmail).Return(emails, nil)

	router.SentEmailsHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []models.Email
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, emails, response)
}

func TestSentEmailsHandler_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router := &EmailRouter{}

	req := httptest.NewRequest(http.MethodGet, "/sent-emails", nil)
	rr := httptest.NewRecorder()

	router.SentEmailsHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestSentEmailsHandler_ErrorGettingEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := &EmailRouter{EmailUseCase: mockEmailUseCase}

	senderEmail := "test@example.com"

	req := httptest.NewRequest(http.MethodGet, "/sent-emails", nil)
	ctx := context.WithValue(req.Context(), "email", senderEmail)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().GetSentEmails(senderEmail).Return(nil, assert.AnError)

	router.SentEmailsHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
