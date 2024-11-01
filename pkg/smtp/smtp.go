package smtp

import (
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPError struct {
	Message string
}

func (e *SMTPError) Error() string {
	return e.Message
}

type SMTPClient struct {
	Host     string
	Port     string
	Username string
	Password string
	Auth     smtp.Auth
}

var sendMail = smtp.SendMail

func NewSMTPClient(host, port, username, password string) *SMTPClient {
	return &SMTPClient{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Auth:     smtp.PlainAuth("", username, password, host),
	}
}

func (c *SMTPClient) SendEmail(from string, to []string, subject, body string) error {
	if len(to) == 0 {
		return fmt.Errorf("список получателей пуст")
	}

	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	msg := formatMessage(from, to, subject, body)

	return sendMail(addr, c.Auth, from, to, []byte(msg))
}

func formatMessage(from string, to []string, subject, body string) string {
	return fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s",
		from,
		strings.Join(to, ", "),
		subject,
		body)
}
