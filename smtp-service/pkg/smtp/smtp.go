package smtp

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type SMTPClient struct {
	Host      string
	Port      string
	Username  string
	Password  string
	UseTLS    bool
	TLSConfig *tls.Config
}

func NewSMTPClient(host, port, username, password string) *SMTPClient {
	return &SMTPClient{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		UseTLS:   true,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         host,
		},
	}
}

func formatMessage(from string, to []string, subject, body string) string {
	msg := "From: " + from + "\r\n" +
		"To: " + strings.Join(to, ", ") + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body
	return msg
}

func (c *SMTPClient) SendEmail(from string, to []string, subject, body string) error {
	if len(to) == 0 {
		return fmt.Errorf("список получателей пуст")
	}

	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
	}

	var conn net.Conn
	var err error

	if c.UseTLS {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, c.TLSConfig)
	} else {
		conn, err = dialer.Dial("tcp", addr)
	}

	if err != nil {
		return fmt.Errorf("ошибка подключения: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	client, err := smtp.NewClient(conn, c.Host)
	if err != nil {
		return fmt.Errorf("ошибка создания SMTP клиента: %v", err)
	}
	defer client.Quit()

	auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("ошибка аутентификации: %v", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("ошибка указания отправителя: %v", err)
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("ошибка указания получателя: %v", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("ошибка начала отправки данных: %v", err)
	}

	msg := formatMessage(from, to, subject, body)
	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("ошибка записи сообщения: %v", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("ошибка завершения отправки данных: %v", err)
	}

	return client.Quit()
}
