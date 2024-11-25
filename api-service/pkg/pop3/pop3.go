package pop3

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"mail/models"
	"net"
	"strings"
	"time"
)

type Pop3Client struct {
	Host      string
	Port      string
	Username  string
	Password  string
	conn      net.Conn
	reader    *bufio.Reader
	writer    *bufio.Writer
	UseTLS    bool
	TLSConfig *tls.Config
}

func NewPop3Client(host, port, username, password string) *Pop3Client {
	return &Pop3Client{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (c *Pop3Client) Connect() error {
	if err := c.dial(); err != nil {
		return err
	}

	if err := c.authenticate(); err != nil {
		return err
	}

	return nil
}

func (c *Pop3Client) dial() error {
	address := fmt.Sprintf("%s:%s", c.Host, c.Port)

	if c.UseTLS {
		tlsConfig := c.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				ServerName: c.Host,
			}
		}

		conn, err := tls.Dial("tcp", address, tlsConfig)
		if err != nil {
			return fmt.Errorf("не удалось установить TLS соединение: %v", err)
		}
		c.conn = conn
	} else {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return fmt.Errorf("не удалось подключиться к POP3 серверу: %v", err)
		}
		c.conn = conn
	}

	c.reader = bufio.NewReader(c.conn)
	c.writer = bufio.NewWriter(c.conn)

	return c.checkResponse("подключение")
}

func (c *Pop3Client) authenticate() error {
	if err := c.sendCommandAndCheck(fmt.Sprintf("USER %s", c.Username), "USER"); err != nil {
		return err
	}

	return c.sendCommandAndCheck(fmt.Sprintf("PASS %s", c.Password), "PASS")
}

func (c *Pop3Client) sendCommandAndCheck(cmd, operation string) error {
	if err := c.sendCommand(cmd); err != nil {
		return err
	}
	return c.checkResponse(operation)
}

func (c *Pop3Client) checkResponse(operation string) error {
	response, err := c.readResponse()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(response, "+OK") {
		return fmt.Errorf("ошибка при %s: %s", operation, response)
	}
	return nil
}

func (c *Pop3Client) sendCommand(cmd string) error {
	_, err := c.writer.WriteString(fmt.Sprintf("%s\r\n", cmd))
	if err != nil {
		return fmt.Errorf("ошибка при отправке команды '%s': %v", cmd, err)
	}
	return c.writer.Flush()
}

func (c *Pop3Client) readResponse() (string, error) {
	response, err := c.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении ответа: %v", err)
	}
	return strings.TrimSpace(response), nil
}

func (c *Pop3Client) readMultilineResponse() (string, error) {
	var lines []string
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("ошибка при чтении многострочного ответа: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "." {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n"), nil
}

func (c *Pop3Client) retrieveEmail(msgNumber int) (string, error) {
	if err := c.sendCommand(fmt.Sprintf("RETR %d", msgNumber)); err != nil {
		return "", err
	}
	response, err := c.readMultilineResponse()
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(response, "+OK") {
		return "", fmt.Errorf("ошибка при выполнении RETR %d: %s", msgNumber, response)
	}
	return response, nil
}

func (c *Pop3Client) FetchEmails(repo models.EmailRepository) error {
	messageCount, err := c.getMessageCount()
	if err != nil {
		return err
	}

	return c.processEmails(messageCount, repo)
}

func (c *Pop3Client) getMessageCount() (int, error) {
	if err := c.sendCommand("STAT"); err != nil {
		return 0, err
	}

	response, err := c.readResponse()
	if err != nil {
		return 0, err
	}

	parts := strings.Split(response, " ")
	if len(parts) < 2 {
		return 0, fmt.Errorf("неверный ответ на STAT: %s", response)
	}

	var count int
	fmt.Sscanf(parts[1], "%d", &count)
	return count, nil
}

func (c *Pop3Client) processEmails(count int, repo models.EmailRepository) error {
	var saveErrors []error
	var savedCount int

	for i := 1; i <= count; i++ {
		if err := c.processSingleEmail(i, repo); err != nil {
			saveErrors = append(saveErrors, fmt.Errorf("письмо #%d: %v", i, err))
			continue
		}
		savedCount++
	}

	if savedCount == 0 && len(saveErrors) > 0 {
		return fmt.Errorf("ошибки при обработке писем: %v", saveErrors)
	}

	return nil
}

func (c *Pop3Client) processSingleEmail(msgNumber int, repo models.EmailRepository) error {
	content, err := c.retrieveEmail(msgNumber)
	if err != nil {
		return err
	}

	email, err := parseEmail(content)
	if err != nil {
		return err
	}

	return repo.SaveEmail(email)
}

func (c *Pop3Client) Quit() error {
	if err := c.sendCommand("QUIT"); err != nil {
		return err
	}
	response, err := c.readResponse()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(response, "+OK") {
		return fmt.Errorf("ошибка при выполнении QUIT: %s", response)
	}
	return c.conn.Close()
}

func parseEmail(raw string) (models.Email, error) {
	email := models.Email{
		IsRead: false,
	}

	lines := strings.Split(raw, "\n")
	bodyIndex := -1

	for i, line := range lines {
		if line == "" {
			bodyIndex = i + 1
			break
		}
		if strings.HasPrefix(line, "From: ") {
			email.Sender_email = strings.TrimSpace(strings.TrimPrefix(line, "From:"))
		} else if strings.HasPrefix(line, "To: ") {
			email.Recipient = strings.TrimSpace(strings.TrimPrefix(line, "To:"))
		} else if strings.HasPrefix(line, "Subject: ") {
			email.Title = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
		} else if strings.HasPrefix(line, "Date: ") {
			dateStr := strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
			parsedDate, err := time.Parse(time.RFC1123Z, dateStr)
			if err != nil {
				return email, fmt.Errorf("не удалось распарсить дату: %v", err)
			}
			email.Sending_date = parsedDate
		}
	}

	if bodyIndex == -1 || bodyIndex >= len(lines) {
		return email, nil
	}
	email.Description = strings.Join(lines[bodyIndex:], "\n")

	return email, nil
}
