package smtp

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"io"
	"math/big"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

type MockSMTPServer struct {
	listener net.Listener
	useTLS   bool
	cert     tls.Certificate
}

func generateTestCert() (tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Corp"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  priv,
	}, nil
}

func NewMockSMTPServer(t *testing.T, useTLS bool) *MockSMTPServer {
	var listener net.Listener
	var err error
	var cert tls.Certificate

	if useTLS {
		cert, err = generateTestCert()
		if err != nil {
			t.Fatalf("Не удалось создать тестовый сертификат: %v", err)
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}

		listener, err = tls.Listen("tcp", "127.0.0.1:0", config)
	} else {
		listener, err = net.Listen("tcp", "127.0.0.1:0")
	}

	if err != nil {
		t.Fatalf("Не удалось создать сервер: %v", err)
	}

	return &MockSMTPServer{
		listener: listener,
		useTLS:   useTLS,
		cert:     cert,
	}
}

func (s *MockSMTPServer) Address() string {
	return s.listener.Addr().String()
}

func (s *MockSMTPServer) Stop() {
	s.listener.Close()
}

func handleSMTPConnection(t *testing.T, conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	if _, err := conn.Write([]byte("220 mock.smtp.server готов\r\n")); err != nil {
		t.Errorf("Ошибка отправки приветствия: %v", err)
		return
	}

	for {
		cmd, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				t.Errorf("Ошибка чтения команды: %v", err)
			}
			return
		}
		cmd = strings.TrimSpace(cmd)
		t.Logf("Сервер получил команду: %s", cmd)

		var response string

		switch {
		case strings.HasPrefix(strings.ToUpper(cmd), "EHLO"):
			response = "250-mock.smtp.server\r\n250 AUTH PLAIN\r\n"
		case strings.HasPrefix(strings.ToUpper(cmd), "AUTH"):
			response = "235 Authentication successful\r\n"
		case strings.HasPrefix(strings.ToUpper(cmd), "MAIL FROM:"):
			response = "250 Отправитель принят\r\n"
		case strings.HasPrefix(strings.ToUpper(cmd), "RCPT TO:"):
			response = "250 Получатель принят\r\n"
		case strings.HasPrefix(strings.ToUpper(cmd), "DATA"):
			response = "354 Начните ввод данных\r\n"
			if _, err := conn.Write([]byte(response)); err != nil {
				t.Errorf("Ошибка отправки ответа DATA: %v", err)
				return
			}
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						t.Errorf("Ошибка чтения данных письма: %v", err)
					}
					return
				}
				line = strings.TrimSpace(line)
				t.Logf("Сервер получил строку данных: %s", line)
				if line == "." {
					break
				}
			}
			response = "250 Сообщение принято\r\n"
		case strings.HasPrefix(strings.ToUpper(cmd), "QUIT"):
			response = "221 До свидания\r\n"
			if _, err := conn.Write([]byte(response)); err != nil {
				t.Errorf("Ошибка отправки ответа QUIT: %v", err)
			}
			return
		default:
			response = "502 Команда не реализована\r\n"
		}

		if _, err := conn.Write([]byte(response)); err != nil {
			t.Errorf("Ошибка отправки ответа: %v", err)
			return
		}
	}
}

func TestSMTPSendEmail(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      []string
		subject string
		body    string
		wantErr bool
	}{
		{
			name:    "Успешная отправка",
			from:    "sender@example.com",
			to:      []string{"recipient@example.com"},
			subject: "Test Subject",
			body:    "Test Body",
			wantErr: false,
		},
		{
			name:    "Пустой получатель",
			from:    "sender@example.com",
			to:      []string{},
			subject: "Test Subject",
			body:    "Test Body",
			wantErr: true,
		},
		{
			name:    "Множественные получатели",
			from:    "sender@example.com",
			to:      []string{"recipient1@example.com", "recipient2@example.com"},
			subject: "Test Subject",
			body:    "Test Body",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr && len(tt.to) == 0 {
				server := NewMockSMTPServer(t, false) // Без TLS
				defer server.Stop()

				go func() {
					conn, err := server.listener.Accept()
					if err != nil {
						return
					}
					handleSMTPConnection(t, conn)
				}()

				addr := server.Address()
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					t.Fatalf("Не удалось разобрать адрес сервера: %v", err)
				}

				client := NewSMTPClient(host, port, "test@example.com", "password")
				client.UseTLS = false

				err = client.SendEmail(tt.from, tt.to, tt.subject, tt.body)
				if err == nil {
					t.Error("SendEmail() должен вернуть ошибку из-за пустого получателя")
				}
				return
			}

			server := NewMockSMTPServer(t, true)
			defer server.Stop()

			var wg sync.WaitGroup
			wg.Add(1)

			// Обработчик соединений
			go func() {
				defer wg.Done()
				conn, err := server.listener.Accept()
				if err != nil {
					t.Errorf("Ошибка принятия соединения: %v", err)
					return
				}
				handleSMTPConnection(t, conn)
			}()

			addr := server.Address()
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				t.Fatalf("Не удалось разобрать адрес сервера: %v", err)
			}

			client := NewSMTPClient(host, port, "test@example.com", "password")
			client.UseTLS = true
			client.TLSConfig = &tls.Config{
				InsecureSkipVerify: true,
			}

			done := make(chan error, 1)
			go func() {
				done <- client.SendEmail(tt.from, tt.to, tt.subject, tt.body)
			}()

			select {
			case err := <-done:
				if (err != nil) != tt.wantErr {
					t.Errorf("SendEmail() error = %v, wantErr %v", err, tt.wantErr)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("Тест превысил время ожидания")
			}

			wg.Wait()
		})
	}
}
