package pop3

import (
	"net/mail"

	"github.com/denisss025/go-pop3-client"
)

type POP3Client struct {
	Host     string
	Port     string
	Username string
	Password string
	Client   *pop3.Client
}

func NewPOP3Client(host, port, username, password string) *POP3Client {
	return &POP3Client{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Client:   &pop3.Client{},
	}
}

func (c *POP3Client) RetrieveMessages(from string) ([]*mail.Message, error) {
	messages, err := c.Client.GetMessages()
	if err != nil {
		return nil, err
	}
	var mailMessages []*mail.Message
	for _, msg := range messages {
		_, mailMsg, err := msg.Retrieve()
		if err != nil {
			return nil, err
		}
		if mailMsg.Header.Get("From") == from {
			mailMessages = append(mailMessages, mailMsg)
		}
	}

	return mailMessages, nil
}
