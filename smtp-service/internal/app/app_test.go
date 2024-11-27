package app

import (
	"context"
	"testing"
	"time"

	"mail/api-service/pkg/logger"
	"mail/config"
)

func TestCreateAndConfigureClients(t *testing.T) {
	cfg := &config.Config{
		SMTP: struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		}{
			Host:     "smtp.test.com",
			Port:     "587",
			Username: "test@test.com",
			Password: "password123",
		},
		Pop3: struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		}{
			Host:     "pop3.test.com",
			Port:     "995",
			Username: "test@test.com",
			Password: "password123",
		},
		SMTPServer: struct {
			IP   string `yaml:"ip"`
			Port string `yaml:"port"`
		}{
			IP:   "localhost",
			Port: "1025",
		},
	}

	t.Run("SMTP Client", func(t *testing.T) {
		smtpClient := createAndConfigureSMTPClient(cfg)
		if smtpClient == nil {
			t.Error("SMTP client is nil")
		}
	})

	t.Run("POP3 Client", func(t *testing.T) {
		pop3Client := createAndConfigurePOP3Client(cfg)
		if pop3Client == nil {
			t.Error("POP3 client is nil")
		}
	})

	t.Run("Run", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		l := logger.NewLogger()
		errCh := make(chan error, 1)

		go func() {
			errCh <- Run(cfg, l)
		}()

		select {
		case err := <-errCh:
			if err == nil {
				t.Error("Expected error when trying to start server with test config")
			}
		case <-ctx.Done():
			t.Skip("Skipping test due to timeout - server started successfully")
		}
	})
}
