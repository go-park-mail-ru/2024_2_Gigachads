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

  err := smtp.SendMail(addr, nil, from, to[0], []byte(msg))


  // // Создаем клиент
  // client, err := smtp.Dial(addr)
  // if err != nil {
  //   return fmt.Errorf("ошибка подключения к SMTP серверу: %v", err)
  // }
  // defer client.Close()

  // //Включаем STARTTLS (обязательно для mail.ru)
  // // tlsConfig := &tls.Config{
  // //   ServerName: c.Host,
  // // }
  // // if err = client.StartTLS(tlsConfig); err != nil {
  // //   return fmt.Errorf("ошибка STARTTLS: %v", err)
  // // }
  // if err = client.Hello(c.Host); err != nil {
// 	return fmt.Errorf("ошибка приветствия SMTP сервера: %v", err)
  // }
  // // Аутентификация (используем учетные данные mail.ru)
  // auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)
  // if err = client.Auth(auth); err != nil {
  //   return fmt.Errorf("ошибка аутентификации: %v", err)
  // }

  // // Отправитель
  // if err = client.Mail(from); err != nil {
  //   return fmt.Errorf("ошибка указания отправителя: %v", err)
  // }

  // // Получатели
  // for _, addr := range to {
  //   if err = client.Rcpt(addr); err != nil {
  //     return fmt.Errorf("ошибка указания получателя: %v", err)
  //   }
  // }

  // // Отправка сообщения
  // w, err := client.Data()
  // if err != nil {
  //   return fmt.Errorf("ошибка начала отправки данных: %v", err)
  // }

  // msg := fmt.Sprintf("From: %s\r\n"+
  //   "To: %s\r\n"+
  //   "Subject: %s\r\n"+
  //   "Content-Type: text/plain; charset=UTF-8\r\n"+
  //   "\r\n"+
  //   "%s", from, to[0], subject, body)

  // _, err = w.Write([]byte(msg))
  // if err != nil {
  //   return fmt.Errorf("ошибка записи сообщения: %v", err)
  // }

  // err = smtp.SendMail(addr, auth, from, to, []byte(msg))
  // if err != nil {
  //   return fmt.Errorf("ошибка отправки сообщения: %v", err)
  // }

  // err = w.Close()
  // if err != nil {
  //   return fmt.Errorf("ошибка завершения отправки: %v", err)
  // }

  // return client.Quit()
}
