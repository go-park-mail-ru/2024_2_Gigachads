package smtp

import (
	"testing"
)

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
