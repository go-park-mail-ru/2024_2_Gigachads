package pop3

import (
	"bufio"
	"fmt"
	"mail/internal/delivery/httpserver/email/mocks"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

type MockPOP3Server struct {
	listener net.Listener
}

func NewMockPOP3Server(t *testing.T) *MockPOP3Server {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Не удалось создать сервер: %v", err)
	}
	return &MockPOP3Server{listener: listener}
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

	server := NewMockPOP3Server(t)
	defer server.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleMockServer(t, server, &wg, emails)

	addr := strings.Split(server.Address(), ":")
	host, port := addr[0], addr[1]

	client := NewPop3Client(host, port, "user", "pass")

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

	server := NewMockPOP3Server(t)
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

	server := NewMockPOP3Server(t)
	defer server.Stop()

	addr := strings.Split(server.Address(), ":")
	host, port := addr[0], addr[1]

	for i, emails := range emailSets {
		var wg sync.WaitGroup
		wg.Add(1)

		go handleMockServer(t, server, &wg, emails)

		client := NewPop3Client(host, port, "user", "pass")

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
