package pop3

import (
  "bufio"
  "fmt"
  "mail/smtp-service/internal/models"
  "net"
  "strings"
  "time"
)

type Pop3Client struct {
  Host     string
  Port     string
  Username string
  Password string
  conn     net.Conn
  reader   *bufio.Reader
  writer   *bufio.Writer
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
  // addr := "pop.mail.ru:995"

  // tlsConfig := &tls.Config{
  //   ServerName:         "pop.mail.ru",
  //   InsecureSkipVerify: true,
  // }

  conn, err := net.Dial("tcp", c.Host+":"+c.Port)
  if err != nil {
    return fmt.Errorf("ошибка подключения к POP3 серверу: %v", err)
  }

  c.conn = conn
  c.reader = bufio.NewReader(conn)
  c.writer = bufio.NewWriter(conn)

  // Устанавливаем таймаут чтения
  c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))

  // Читаем приветствие сервера
  response, err := c.readResponse()
  if err != nil {
    c.conn.Close()
    return fmt.Errorf("ошибка чтения приветствия: %v", err)
  }
  fmt.Printf("Приветствие сервера: %s\n", response)

  // Отправляем USER
  if err := c.sendCommand(fmt.Sprintf("USER %s", c.Username)); err != nil {
    c.conn.Close()
    return fmt.Errorf("ошибка отправки USER: %v", err)
  }
  response, err = c.readResponse()
  if err != nil {
    c.conn.Close()
    return fmt.Errorf("ошибка ответа на USER: %v", err)
  }
  fmt.Printf("Ответ на USER: %s\n", response)

  // Отправляем PASS
  if err := c.sendCommand(fmt.Sprintf("PASS %s", c.Password)); err != nil {
    c.conn.Close()
    return fmt.Errorf("ошибка отправки PASS: %v", err)
  }
  response, err = c.readResponse()
  if err != nil {
    c.conn.Close()
    return fmt.Errorf("ошибка ответа на PASS: %v", err)
  }
  fmt.Printf("Ответ на PASS: %s\n", response)

  return nil
}

func (c *Pop3Client) FetchEmails(repo models.EmailRepositorySMTP) error {
  if c.conn == nil {
    return fmt.Errorf("соединение не установлено")
  }

  c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))

  if err := c.sendCommand("STAT"); err != nil {
    return fmt.Errorf("ошибка отправки STAT: %v", err)
  }

  response, err := c.readResponse()
  if err != nil {
    return fmt.Errorf("ошибка чтения ответа STAT: %v", err)
  }
  fmt.Printf("Ответ на STAT: %s\n", response)

  var count, size int
  _, err = fmt.Sscanf(response, "+OK %d %d", &count, &size)
  if err != nil {
    return fmt.Errorf("ошибка парсинга STAT (%s): %v", response, err)
  }

  fmt.Printf("Найдено писем: %d, общий размер: %d байт\n", count, size)

  for i := 1; i <= count; i++ {
    // Обновляем таймаут для каждого письма
    c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))

    email, err := c.retrieveEmail(i)
    if err != nil {
      fmt.Printf("Ошибка получения письма %d: %v\n", i, err)
      continue
    }

    parsedEmail, err := parseEmail(email)
    if err != nil {
      fmt.Printf("Ошибка парсинга письма %d: %v\n", i, err)
      continue
    }

    err = repo.SaveEmail(parsedEmail)
    if err != nil {
      fmt.Printf("Ошибка сохранения письма %d: %v\n", i, err)
      continue
    }
  }

  return nil
}

func (c *Pop3Client) retrieveEmail(msgNum int) (string, error) {
  err := c.sendCommand(fmt.Sprintf("RETR %d", msgNum))
  if err != nil {
    return "", fmt.Errorf("ошибка отправки команды RETR: %v", err)
  }

  response, err := c.readResponse()
  if err != nil {
    return "", fmt.Errorf("ошибка чтения ответа на RETR: %v", err)
  }
  fmt.Printf("Ответ на RETR %d: %s\n", msgNum, response)

  var builder strings.Builder
  for {
    c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
    
    line, err := c.reader.ReadString('\n')
    if err != nil {
      return "", fmt.Errorf("ошибка чтения строки письма: %v", err)
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
  fmt.Printf("Отправка команды: %s\n", command)
  _, err := fmt.Fprintf(c.conn, "%s\r\n", command)
  if err != nil {
    return fmt.Errorf("ошибка отправки команды: %v", err)
  }
  err = c.writer.Flush()
  if err != nil {
    return fmt.Errorf("ошибка flush: %v", err)
  }
  return nil
}

func (c *Pop3Client) readResponse() (string, error) {
  response, err := c.reader.ReadString('\n')
  if err != nil {
    return "", fmt.Errorf("ошибка чтения ответа: %v", err)
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
  if c.conn != nil {
    if err := c.sendCommand("QUIT"); err != nil {
      c.conn.Close()
      return fmt.Errorf("ошибка отправки команды QUIT: %v", err)
    }
    return c.conn.Close()
  }
  return nil
}
