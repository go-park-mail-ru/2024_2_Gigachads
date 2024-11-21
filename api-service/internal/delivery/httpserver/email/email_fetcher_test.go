package email

import (
	"context"
	"mail/internal/delivery/httpserver/email/mocks"
	"mail/internal/models"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestEmailFetcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)

	mockEmailUseCase.EXPECT().
		FetchEmailsViaPOP3().
		Return(nil).
		AnyTimes()

	mockEmailUseCase.EXPECT().
		Inbox(gomock.Any()).
		Return([]models.Email{}, nil).
		AnyTimes()

	mockEmailUseCase.EXPECT().
		SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	mockEmailUseCase.EXPECT().
		ForwardEmail(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	mockEmailUseCase.EXPECT().
		ReplyEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	mockEmailUseCase.EXPECT().
		GetEmailByID(gomock.Any()).
		Return(models.Email{}, nil).
		AnyTimes()

	mockEmailUseCase.EXPECT().
		ChangeStatus(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	mockEmailUseCase.EXPECT().
		GetSentEmails(gomock.Any()).
		Return([]models.Email{}, nil).
		AnyTimes()

	mockEmailUseCase.EXPECT().
		SaveEmail(gomock.Any()).
		Return(nil).
		AnyTimes()

	mockEmailUseCase.EXPECT().
		DeleteEmails(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	fetcher := NewEmailFetcher(mockEmailUseCase)

	testTicker := make(chan time.Time)
	originalTicker := NewTicker
	NewTicker = func(d time.Duration) *time.Ticker {
		return &time.Ticker{
			C: testTicker,
		}
	}
	defer func() {
		NewTicker = originalTicker
	}()

	fetcher.Start()
	testTicker <- time.Now()
	time.Sleep(100 * time.Millisecond)
}

func TestEmailFetcherNilService(t *testing.T) {
	fetcher := NewEmailFetcher(nil)
	fetcher.Start()
}

func TestSessionOperations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSession := mocks.NewMockSessionRepository(ctrl)
	ctx := context.Background()

	expectedSession := &models.Session{
		Name:      "session",
		ID:        "test-hash",
		Time:      time.Now().Add(24 * time.Hour),
		UserLogin: "test@mail.com",
	}

	mockSession.EXPECT().
		CreateSession(ctx, "test@mail.com").
		Return(expectedSession, nil).
		Times(1)

	mockSession.EXPECT().
		GetSession(ctx, "test-hash").
		Return("test@mail.com", nil).
		Times(1)

	mockSession.EXPECT().
		DeleteSession(ctx, "test-hash").
		Return(nil).
		Times(1)

	session, err := mockSession.CreateSession(ctx, "test@mail.com")
	if err != nil {
		t.Errorf("CreateSession failed: %v", err)
	}
	if session.ID != expectedSession.ID {
		t.Errorf("Expected session ID %s, got %s", expectedSession.ID, session.ID)
	}

	userEmail, err := mockSession.GetSession(ctx, "test-hash")
	if err != nil {
		t.Errorf("GetSession failed: %v", err)
	}
	if userEmail != "test@mail.com" {
		t.Errorf("Expected email test@mail.com, got %s", userEmail)
	}

	err = mockSession.DeleteSession(ctx, "test-hash")
	if err != nil {
		t.Errorf("DeleteSession failed: %v", err)
	}
}
