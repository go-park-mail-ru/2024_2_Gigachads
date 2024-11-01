package email

import (
	"testing"
	"time"

	"mail/internal/delivery/httpserver/email/mocks"

	"github.com/golang/mock/gomock"
)

func TestEmailFetcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailService := mocks.NewMockEmailUseCase(ctrl)

	mockEmailService.EXPECT().
		FetchEmailsViaPOP3().
		Return(nil).
		MinTimes(2)

	originalTicker := newTicker
	defer func() { newTicker = originalTicker }()

	done := make(chan bool)

	newTicker = func(d time.Duration) *time.Ticker {
		return time.NewTicker(100 * time.Millisecond)
	}

	fetcher := NewEmailFetcher(mockEmailService)
	fetcher.Start()

	go func() {
		time.Sleep(350 * time.Millisecond)
		done <- true
	}()

	<-done
}
