package usecase

import (
	"fmt"
	"time"

	models "mail/internal/models"
)

type EmailService struct {
	EmailRepo   models.EmailRepository
	SessionRepo models.SessionRepository
	SMTPRepo    models.SMTPRepository
}

func NewEmailService(erepo models.EmailRepository, srepo models.SessionRepository, smtprepo models.SMTPRepository) *EmailService {
	return &EmailService{
		EmailRepo:   erepo,
		SessionRepo: srepo,
		SMTPRepo:    smtprepo,
	}
}

func (es *EmailService) Inbox(sessionID string) ([]models.Email, error) {
	session, err := es.SessionRepo.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	return es.EmailRepo.Inbox(session.UserLogin)
}

func (es *EmailService) SendEmail(from string, to []string, subject string, body string) error {
	return es.SMTPRepo.SendEmail(from, to, subject, body)
}

func (es *EmailService) ForwardEmail(from string, to []string, originalEmail models.Email) error {
	forwardSubject := "Fwd: " + originalEmail.Title
	forwardBody := fmt.Sprintf(`
---------- Forwarded message ---------
From: %s
Date: %s
Subject: %s

%s
`, originalEmail.Sender_email, originalEmail.Sending_date.Format(time.RFC1123),
		originalEmail.Title, originalEmail.Description)

	return es.SMTPRepo.SendEmail(from, to, forwardSubject, forwardBody)
}

func (es *EmailService) ReplyEmail(from string, to string, originalEmail models.Email, replyText string) error {
	replySubject := "Re: " + originalEmail.Title
	replyBody := fmt.Sprintf(`
%s

On %s, %s wrote:
> %s
`, replyText, originalEmail.Sending_date.Format(time.RFC1123),
		originalEmail.Sender_email, originalEmail.Description)

	return es.SMTPRepo.SendEmail(from, []string{to}, replySubject, replyBody)
}

func (es *EmailService) GetEmailByID(id int) (models.Email, error) {
	return es.EmailRepo.GetEmailByID(id)
}
