package pop3

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"mail/smtp-service/internal/models"
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
        UseTLS:   true,
        TLSConfig: &tls.Config{
            ServerName:         host,
            InsecureSkipVerify: true,
            MinVersion:         tls.VersionTLS12,
        },
    }
}

func (c *Pop3Client) Connect() error {
    address := fmt.Sprintf("%s:%s", c.Host, c.Port)
    
    dialer := &net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
    }

    var err error
    if c.UseTLS {
        c.conn, err = tls.DialWithDialer(dialer, "tcp", address, c.TLSConfig)
    } else {
        c.conn, err = dialer.Dial("tcp", address)
    }

    if err != nil {
        return fmt.Errorf("ошибка подключения к POP3 серверу: %v", err)
    }

    c.reader = bufio.NewReader(c.conn)
    c.writer = bufio.NewWriter(c.conn)

    _, err = c.readResponse()
    if err != nil {
        return err
    }

    if err := c.sendCommand(fmt.Sprintf("USER %s", c.Username)); err != nil {
        return err
    }
    if err := c.sendCommand(fmt.Sprintf("PASS %s", c.Password)); err != nil {
        return err
    }

    return nil
}

func (c *Pop3Client) FetchEmails(repo models.EmailRepositorySMTP) error {
    // Получаем количество писем
    err := c.sendCommand("STAT")
    if err != nil {
        return err
    }
    response, err := c.readResponse()
    if err != nil {
        return err
    }

    var count int
    _, err = fmt.Sscanf(response, "+OK %d", &count)
    if err != nil {
        return fmt.Errorf("ошибка парсинга STAT: %v", err)
    }

    for i := 1; i <= count; i++ {
        email, err := c.retrieveEmail(i)
        if err != nil {
            return err
        }

        parsedEmail, err := parseEmail(email)
        if err != nil {
            continue 
        }

        err = repo.SaveEmail(parsedEmail)
        if err != nil {
            return err
        }
    }

    return nil
}

func (c *Pop3Client) retrieveEmail(msgNum int) (string, error) {
    err := c.sendCommand(fmt.Sprintf("RETR %d", msgNum))
    if err != nil {
        return "", err
    }

    var builder strings.Builder
    for {
        line, err := c.reader.ReadString('\n')
        if err != nil {
            return "", err
        }

        line = strings.TrimRight(line, "\r\n")
        if line == "." {
            break
        }
        if strings.HasPrefix(line, "..") {
            line = line[1:]
        }
        builder.WriteString(line)
        builder.WriteString("\n")
    }

    return builder.String(), nil
}

func (c *Pop3Client) sendCommand(command string) error {
    _, err := fmt.Fprintf(c.conn, "%s\r\n", command)
    if err != nil {
        return err
    }
    return c.writer.Flush()
}

func (c *Pop3Client) readResponse() (string, error) {
    response, err := c.reader.ReadString('\n')
    if err != nil {
        return "", err
    }
    response = strings.TrimRight(response, "\r\n")

    if !strings.HasPrefix(response, "+OK") {
        return "", fmt.Errorf("ошибка сервера: %s", response)
    }
    return response, nil
}

func parseEmail(raw string) (models.Email, error) {
    email := models.Email{
        IsRead:       false,
        Sending_date: time.Now(),
    }

    headers := strings.Split(raw, "\n\n")[0]
    for _, line := range strings.Split(headers, "\n") {
        if strings.HasPrefix(line, "From:") {
            email.Sender_email = strings.TrimSpace(strings.TrimPrefix(line, "From:"))
        } else if strings.HasPrefix(line, "To:") {
            email.Recipient = strings.TrimSpace(strings.TrimPrefix(line, "To:"))
        } else if strings.HasPrefix(line, "Subject:") {
            email.Title = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
        } else if strings.HasPrefix(line, "Date:") {
            dateStr := strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
            if date, err := time.Parse(time.RFC1123Z, dateStr); err == nil {
                email.Sending_date = date
            }
        }
    }

    parts := strings.SplitN(raw, "\n\n", 2)
    if len(parts) > 1 {
        email.Description = strings.TrimSpace(parts[1])
    }

    return email, nil
}

func (c *Pop3Client) Quit() error {
	return nil
}
