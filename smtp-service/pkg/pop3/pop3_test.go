package pop3

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"mail/smtp-service/pkg/mocks"

	"fmt"
	"math/big"
	"net"
	"strings"

	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

type MockPOP3Server struct {
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
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24),

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

func NewMockPOP3Server(t *testing.T, useTLS bool) *MockPOP3Server {
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
		}

		listener, err = tls.Listen("tcp", "127.0.0.1:0", config)
	} else {
		listener, err = net.Listen("tcp", "127.0.0.1:0")
	}

	if err != nil {
		t.Fatalf("Не удалось создать сервер: %v", err)
	}

	return &MockPOP3Server{
		listener: listener,
		useTLS:   useTLS,
		cert:     cert,
	}
}

func (s *MockPOP3Server) Address() string {
	return s.listener.Addr().String()
}

func (s *MockPOP3Server) Stop() {
	s.listener.Close()
}

func handleMockServer(t *testing.T, server *MockPOP3Server, wg *sync.WaitGroup, emails []string) {
	defer wg.Done()
	conn, err := server.listener.Accept()
	if err != nil {
		if !strings.Contains(err.Error(), "use of closed network connection") {
			t.Errorf("Ошибка при принятии соединения: %v", err)
		}
		return
	}
	defer conn.Close()

	conn.Write([]byte("+OK POP3 сервер готов\r\n"))
	reader := bufio.NewReader(conn)
	emailIndex := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Ошибка чтения: %v", err)
			}
			return
		}

		cmd := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(cmd, "USER"), strings.HasPrefix(cmd, "PASS"):
			conn.Write([]byte("+OK\r\n"))
		case cmd == "STAT":
			conn.Write([]byte(fmt.Sprintf("+OK %d 200\r\n", len(emails))))
		case strings.HasPrefix(cmd, "RETR"):
			if emailIndex < len(emails) {
				conn.Write([]byte(emails[emailIndex]))
				emailIndex++
			} else {
				conn.Write([]byte("-ERR Нет такого сообщения\r\n"))
			}
		case cmd == "QUIT":
			conn.Write([]byte("+OK До свидания\r\n"))
			return
		default:
			conn.Write([]byte("-ERR Неизвестная команда\r\n"))
		}
	}
}

func TestPOP3MultipleEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmailRepository(ctrl)
	repo.EXPECT().SaveEmail(gomock.Any()).Return(nil).Times(2)

	emails := []string{
		"+OK\r\n" +
			"From: sender1@example.com\r\n" +
			"To: recipient1@example.com\r\n" +
			"Subject: Первое письмо\r\n" +
			"Date: Mon, 15 Jan 2024 10:00:00 +0000\r\n" +
			"\r\n" +
			"Тело первого сообщения\r\n" +
			".\r\n",

		"+OK\r\n" +
			"From: sender2@example.com\r\n" +
			"To: recipient2@example.com\r\n" +
			"Subject: Второе письмо\r\n" +
			"Date: Mon, 15 Jan 2024 11:00:00 +0000\r\n" +
			"\r\n" +
			"Тело второго сообщения\r\n" +
			".\r\n",
	}

	server := NewMockPOP3Server(t, false)
	defer server.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleMockServer(t, server, &wg, emails)

	addr := strings.Split(server.Address(), ":")
	host, port := addr[0], addr[1]

	client := NewPop3Client(host, port, "user", "pass")
	client.UseTLS = false

	err := client.Connect()
	if err != nil {
		t.Fatalf("Ошибка подключения: %v", err)
	}
	defer client.Quit()

	if client.conn != nil {
		client.conn.SetDeadline(time.Now().Add(5 * time.Second))
	}

	err = client.FetchEmails(repo)
	if err != nil {
		t.Fatalf("Ошибка при получении писем: %v", err)
	}

	err = client.Quit()
	if err != nil {
		t.Fatalf("Ошибка при закрытии соединения: %v", err)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Таймаут при ожидании завершения сервера")
	}
}

func TestPOP3ConcurrentFetching(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	emails := []string{
		"+OK\r\n" +
			"From: sender1@example.com\r\n" +
			"To: recipient1@example.com\r\n" +
			"Subject: Письмо 1\r\n" +
			"Date: Mon, 15 Jan 2024 10:00:00 +0000\r\n" +
			"\r\n" +
			"Тело сообщения 1\r\n" +
			".\r\n",
	}

	server := NewMockPOP3Server(t, false)
	defer server.Stop()

	done := make(chan bool)

	go func() {
		defer close(done)
		successfulConnections := 0
		for successfulConnections < 3 {
			conn, err := server.listener.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					t.Errorf("Ошибка при принятии соединения: %v", err)
				}
				return
			}
			successfulConnections++

			go func(conn net.Conn) {
				defer conn.Close()

				conn.Write([]byte("+OK POP3 сервер готов\r\n"))

				reader := bufio.NewReader(conn)
				emailIndex := 0

				for {
					line, err := reader.ReadString('\n')
					if err != nil {
						return
					}

					cmd := strings.TrimSpace(line)
					switch {
					case strings.HasPrefix(cmd, "USER"), strings.HasPrefix(cmd, "PASS"):
						conn.Write([]byte("+OK\r\n"))
					case cmd == "STAT":
						conn.Write([]byte(fmt.Sprintf("+OK %d 200\r\n", len(emails))))
					case strings.HasPrefix(cmd, "RETR"):
						if emailIndex < len(emails) {
							conn.Write([]byte(emails[emailIndex]))
							emailIndex++
						} else {
							conn.Write([]byte("-ERR Нет такого сообщения\r\n"))
						}
					case cmd == "QUIT":
						conn.Write([]byte("+OK До свидания\r\n"))
						return
					default:
						conn.Write([]byte("-ERR Неизвестная команда\r\n"))
					}
				}
			}(conn)
		}
	}()

	addr := strings.Split(server.Address(), ":")
	host, port := addr[0], addr[1]

	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go func(clientNum int) {
			defer wg.Done()

			repo := mocks.NewMockEmailRepository(ctrl)
			repo.EXPECT().SaveEmail(gomock.Any()).Return(nil).Times(1)

			client := NewPop3Client(host, port, "user", "pass")
			client.UseTLS = false

			timeout := time.After(5 * time.Second)
			connected := make(chan error, 1)
			go func() {
				connected <- client.Connect()
			}()

			select {
			case err := <-connected:
				if err != nil {
					t.Errorf("Клиент %d: ошибка подключения: %v", clientNum, err)
					return
				}
			case <-timeout:
				t.Errorf("Клиент %d: таймаут при подключении", clientNum)
				return
			}

			defer client.Quit()

			if client.conn != nil {
				client.conn.SetDeadline(time.Now().Add(5 * time.Second))
			}

			err := client.FetchEmails(repo)
			if err != nil {
				t.Errorf("Клиент %d: ошибка при получении писем: %v", clientNum, err)
				return
			}
		}(i)
	}

	wg.Wait()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Таймаут при ожидании завершения сервера")
	}
}

func TestPOP3SequentialFetching(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmailRepository(ctrl)
	repo.EXPECT().SaveEmail(gomock.Any()).Return(nil).Times(3)

	emailSets := [][]string{
		{
			"+OK\r\n" +
				"From: sender1@example.com\r\n" +
				"To: recipient1@example.com\r\n" +
				"Subject: Первый пакет\r\n" +
				"Date: Mon, 15 Jan 2024 10:00:00 +0000\r\n" +
				"\r\n" +
				"Первое сообщение\r\n" +
				".\r\n",
		},
		{
			"+OK\r\n" +
				"From: sender2@example.com\r\n" +
				"To: recipient2@example.com\r\n" +
				"Subject: Второй пакет\r\n" +
				"Date: Mon, 15 Jan 2024 11:00:00 +0000\r\n" +
				"\r\n" +
				"Второе сообщение\r\n" +
				".\r\n",
		},
		{
			"+OK\r\n" +
				"From: sender3@example.com\r\n" +
				"To: recipient3@example.com\r\n" +
				"Subject: Третий пакет\r\n" +
				"Date: Mon, 15 Jan 2024 12:00:00 +0000\r\n" +
				"\r\n" +
				"Третье сообщение\r\n" +
				".\r\n",
		},
	}

	server := NewMockPOP3Server(t, false)
	defer server.Stop()

	addr := strings.Split(server.Address(), ":")
	host, port := addr[0], addr[1]

	for i, emails := range emailSets {
		var wg sync.WaitGroup
		wg.Add(1)

		go handleMockServer(t, server, &wg, emails)

		client := NewPop3Client(host, port, "user", "pass")
		client.UseTLS = false

		err := client.Connect()
		if err != nil {
			t.Fatalf("Пакет %d: ошибка подключения: %v", i+1, err)
		}

		if client.conn != nil {
			client.conn.SetDeadline(time.Now().Add(5 * time.Second))
		}

		err = client.FetchEmails(repo)
		if err != nil {
			t.Fatalf("Пакет %d: ошибка при получении писем: %v", i+1, err)
		}

		err = client.Quit()
		if err != nil {
			t.Fatalf("Пакет %d: ошибка при закрытии соединения: %v", i+1, err)
		}

		time.Sleep(100 * time.Millisecond)
	}

}

func TestPOP3WithTLS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmailRepository(ctrl)
	repo.EXPECT().SaveEmail(gomock.Any()).Return(nil).Times(1)

	emails := []string{
		"+OK\r\n" +
			"From: sender1@example.com\r\n" +
			"To: recipient1@example.com\r\n" +
			"Subject: TLS Test\r\n" +
			"Date: Mon, 15 Jan 2024 10:00:00 +0000\r\n" +
			"\r\n" +
			"Тестовое TLS сообщение\r\n" +
			".\r\n",
	}

	server := NewMockPOP3Server(t, true)
	defer server.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleMockServer(t, server, &wg, emails)

	addr := strings.Split(server.Address(), ":")
	host, port := addr[0], addr[1]

	client := NewPop3Client(host, port, "user", "pass")
	client.UseTLS = true

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	client.TLSConfig = tlsConfig

	err := client.Connect()
	if err != nil {
		t.Fatalf("Ошибка подключения: %v", err)
	}

	if client.conn != nil {
		client.conn.SetDeadline(time.Now().Add(5 * time.Second))
	}

	err = client.FetchEmails(repo)
	if err != nil {
		t.Fatalf("Ошибка при получении писем через TLS: %v", err)
	}

	err = client.Quit()
	if err != nil {
		t.Fatalf("Ошибка при закрытии TLS соединения: %v", err)
	}

	wg.Wait()
}
