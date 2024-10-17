package smtp

import (
	"net/smtp"
	"strings"
)

type SMTPClient struct {
	Host     string
	Port     string
	Username string
	Password string
	Auth     smtp.Auth
}

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
	msg := "From: " + from + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	addr := c.Host + ":" + c.Port
	return smtp.SendMail(addr, c.Auth, from, to, []byte(msg))
}
