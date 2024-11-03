package smtp

import (
	"fmt"
	"net/smtp"
	"strings"
	"testing"
)

func TestSMTPError(t *testing.T) {
	err := &SMTPError{Message: "test error"}
	if err.Error() != "test error" {
		t.Errorf("SMTPError.Error() = %v, want %v", err.Error(), "test error")
	}
}

func TestNewSMTPClient(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     string
		username string
		password string
	}{
		{
			name:     "Создание клиента с валидными данными",
			host:     "smtp.example.com",
			port:     "587",
			username: "test@example.com",
			password: "password123",
		},
		{
			name:     "Создание клиента с пустым портом",
			host:     "smtp.example.com",
			port:     "",
			username: "test@example.com",
			password: "password123",
		},
		{
			name:     "Создание клиента с пустыми данными",
			host:     "",
			port:     "",
			username: "",
			password: "",
		},
		{
			name:     "Создание клиента с специальными символами",
			host:     "smtp.test-server.com",
			port:     "465",
			username: "test.user+label@example.com",
			password: "pass!@#$%^&*()",
		},
		{
			name:     "Создание клиента с нестандартным портом",
			host:     "smtp.custom.com",
			port:     "2525",
			username: "admin@custom.com",
			password: "adminpass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewSMTPClient(tt.host, tt.port, tt.username, tt.password)

			if client.Host != tt.host {
				t.Errorf("неверный хост: получили %v, ожидали %v", client.Host, tt.host)
			}
			if client.Port != tt.port {
				t.Errorf("неверный порт: получили %v, ожидали %v", client.Port, tt.port)
			}
			if client.Username != tt.username {
				t.Errorf("неверное имя пользователя: получили %v, ожидали %v", client.Username, tt.username)
			}
			if client.Password != tt.password {
				t.Errorf("неверный пароль: получили %v, ожидали %v", client.Password, tt.password)
			}
			if client.Auth == nil {
				t.Error("Auth не должен быть nil")
			}
		})
	}
}

func TestSendEmail(t *testing.T) {
	tests := []struct {
		name        string
		from        string
		to          []string
		subject     string
		body        string
		mockSendErr error
		wantErr     bool
	}{
		{
			name:    "Успешная отправка письма",
			from:    "sender@example.com",
			to:      []string{"recipient@example.com"},
			subject: "Тестовое письмо",
			body:    "Тело письма",
			wantErr: false,
		},
		{
			name:        "Ошибка при отправке",
			from:        "sender@example.com",
			to:          []string{"recipient@example.com"},
			subject:     "Тестовое письмо",
			body:        "Тело письма",
			mockSendErr: fmt.Errorf("ошибка отправки"),
			wantErr:     true,
		},
		{
			name:    "Пустой получатель",
			from:    "sender@example.com",
			to:      []string{},
			subject: "Тестовое письмо",
			body:    "Тело письма",
			wantErr: true,
		},
		{
			name:    "Множественные получатели",
			from:    "sender@example.com",
			to:      []string{"recipient1@example.com", "recipient2@example.com"},
			subject: "Тестовое письмо",
			body:    "Тело письма",
			wantErr: false,
		},
		{
			name:    "Длинное тело письма",
			from:    "sender@example.com",
			to:      []string{"recipient@example.com"},
			subject: "Длинное письмо",
			body:    strings.Repeat("Long content ", 1000),
			wantErr: false,
		},
		{
			name:    "Специальные символы в теме",
			from:    "sender@example.com",
			to:      []string{"recipient@example.com"},
			subject: "Test !@#$%^&*()_+ Тест",
			body:    "Тело письма",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &SMTPClient{
				Host:     "smtp.example.com",
				Port:     "587",
				Username: "test@example.com",
				Password: "password123",
				Auth:     smtp.PlainAuth("", "test@example.com", "password123", "smtp.example.com"),
			}

			originalSendMail := sendMail
			defer func() { sendMail = originalSendMail }()

			sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
				return tt.mockSendErr
			}

			err := client.SendEmail(tt.from, tt.to, tt.subject, tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormatMessage(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      []string
		subject string
		body    string
		want    string
	}{
		{
			name:    "Базовое форматирование",
			from:    "sender@example.com",
			to:      []string{"recipient@example.com"},
			subject: "Тест",
			body:    "Тестовое сообщение",
			want: "From: sender@example.com\r\n" +
				"To: recipient@example.com\r\n" +
				"Subject: Тест\r\n" +
				"Content-Type: text/plain; charset=UTF-8\r\n" +
				"\r\n" +
				"Тестовое сообщение",
		},
		{
			name:    "Множественные получатели",
			from:    "sender@example.com",
			to:      []string{"recipient1@example.com", "recipient2@example.com"},
			subject: "Тест",
			body:    "Тестовое сообщение",
			want: "From: sender@example.com\r\n" +
				"To: recipient1@example.com, recipient2@example.com\r\n" +
				"Subject: Тест\r\n" +
				"Content-Type: text/plain; charset=UTF-8\r\n" +
				"\r\n" +
				"Тестовое сообщение",
		},
		{
			name:    "Пустое тело письма",
			from:    "sender@example.com",
			to:      []string{"recipient@example.com"},
			subject: "Пустое письмо",
			body:    "",
			want: "From: sender@example.com\r\n" +
				"To: recipient@example.com\r\n" +
				"Subject: Пустое письмо\r\n" +
				"Content-Type: text/plain; charset=UTF-8\r\n" +
				"\r\n",
		},
		{
			name:    "Специальные символы",
			from:    "test.user+label@example.com",
			to:      []string{"special!user@test.com"},
			subject: "Test !@#$%^&*()_+ Тест",
			body:    "Body with спецсимволы !@#$%^&*()",
			want: "From: test.user+label@example.com\r\n" +
				"To: special!user@test.com\r\n" +
				"Subject: Test !@#$%^&*()_+ Тест\r\n" +
				"Content-Type: text/plain; charset=UTF-8\r\n" +
				"\r\n" +
				"Body with спецсимволы !@#$%^&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatMessage(tt.from, tt.to, tt.subject, tt.body)
			if got != tt.want {
				t.Errorf("formatMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
