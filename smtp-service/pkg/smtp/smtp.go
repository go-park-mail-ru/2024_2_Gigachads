package smtp

import (
  "fmt"
  "net/smtp"
  //"crypto/tls"
)

type SMTPClient struct {
  Host     string
  Port     string
  Username string
  Password string
}

func NewSMTPClient(host, port, username, password string) *SMTPClient {
  return &SMTPClient{
    Host:     host,
    Port:     port,
    Username: username,
    Password: password,
  }
}

func (c *SMTPClient) SendEmail(from string, to []string, subject, body string) error {
  if len(to) == 0 {
    return fmt.Errorf("список получателей пуст")
  }

  msg := fmt.Sprintf("From: %s\r\n"+
    "To: %s\r\n"+
    "Subject: %s\r\n"+
    "Content-Type: text/plain; charset=UTF-8\r\n"+
    "\r\n"+
    "%s", from, to[0], subject, body)

  addr := c.Host+":"+c.Port

  err := smtp.SendMail(addr, nil, from, to, []byte(msg))
  if err != nil {
    return fmt.Errorf("ошибка отправки сообщения: %v", err)
  }
  return nil
}
