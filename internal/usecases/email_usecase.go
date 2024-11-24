package usecase

import (
	"fmt"
	models "mail/internal/models"
	"time"
)

type EmailService struct {
	EmailRepo   models.EmailRepository
	SessionRepo models.SessionRepository
	SMTPRepo    models.SMTPRepository
	POP3Repo    models.POP3Repository
}

func NewEmailService(
	erepo models.EmailRepository,
	srepo models.SessionRepository,
	smtprepo models.SMTPRepository,
	pop3repo models.POP3Repository,
) *EmailService {
	return &EmailService{
		EmailRepo:   erepo,
		SessionRepo: srepo,
		SMTPRepo:    smtprepo,
		POP3Repo:    pop3repo,
	}
}

func (es *EmailService) Inbox(email string) ([]models.Email, error) {
	// return es.EmailRepo.Inbox(email)
	return es.EmailRepo.GetFolderEmails(email, "Входящие")
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

func (es *EmailService) FetchEmailsViaPOP3() error {
	if err := es.POP3Repo.Connect(); err != nil {
		return err
	}
	defer es.POP3Repo.Quit()

	if err := es.POP3Repo.FetchEmails(es.EmailRepo); err != nil {
		return err
	}

	return nil
}

func (es *EmailService) GetSentEmails(email string) ([]models.Email, error) {
	//return es.EmailRepo.GetSentEmails(senderEmail)
	return es.EmailRepo.GetFolderEmails(email, "Отправленные")
}

func (s *EmailService) SaveEmail(email models.Email) error {
	return s.EmailRepo.SaveEmail(email)
}

func (es *EmailService) ChangeStatus(id int, status bool) error {
	return es.EmailRepo.ChangeStatus(id, status)
}

func (es *EmailService) DeleteEmails(userEmail string, messageIDs []int) error {
	for _, elem := range messageIDs {
		folder, err := es.EmailRepo.GetMessageFolder(elem)
		if err != nil {
			return err
		}
		if folder == "Корзина" {
			return es.EmailRepo.DeleteEmails(userEmail, messageIDs)
		} else {
			err = es.EmailRepo.ChangeEmailFolder(elem, userEmail, "Корзина")
			if err != nil {
				return err
			}
		}
	}
	return nil
}


func (es *EmailService) GetFolders(email string) ([]string, error) {
	emails, err := es.EmailRepo.GetFolders(email)
	if err != nil {
		return nil, err
	}
	if len(emails) == 0 {
		es.EmailRepo.CreateFolder(email, "Входящие")
		es.EmailRepo.CreateFolder(email, "Отправленные")
		es.EmailRepo.CreateFolder(email, "Спам")
		es.EmailRepo.CreateFolder(email, "Черновики")
		es.EmailRepo.CreateFolder(email, "Корзина")
	}
	return es.EmailRepo.GetFolders(email)
}

func (es *EmailService) GetFolderEmails(email string, folderName string) ([]models.Email, error) {
	emails, err := es.EmailRepo.GetFolderEmails(email, folderName)
	if err != nil {
		return nil, err
	}
	if folderName == "Отправленные" {
		for i, _ := range emails {
			temp := emails[i].Sender_email
			emails[i].Sender_email = emails[i].Recipient
			emails[i].Recipient = temp
		}
	}

	// for i, _ := range emails {
	// 	if emails[i].Sender_email == email {
	// 		temp := emails[i].Sender_email
	// 		emails[i].Sender_email = emails[i].Recipient
	// 		emails[i].Recipient = temp
	// 	}
	// }
	
	return emails, nil
}

func (es *EmailService) CreateFolder(email string, folderName string) error {
	if ok, err := es.EmailRepo.CheckFolder(email, folderName); ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("folder_already_exists")	
	}
	return es.EmailRepo.CreateFolder(email, folderName)
}

func (es *EmailService) DeleteFolder(email string, folderName string) error {
	if folderName == "Входящие" || folderName == "Отправленные" || folderName == "Спам" || folderName == "Черновики" || folderName == "Корзина" {
		return fmt.Errorf("unable_to_delete_folder")	
	}
	return es.EmailRepo.DeleteFolder(email, folderName)
}

func (es *EmailService) RenameFolder(email string, folderName string, newFolderName string) error {
	if folderName == "Входящие" || folderName == "Отправленные" || folderName == "Спам" || folderName == "Черновики" || folderName == "Корзина" {
		return fmt.Errorf("unable_to_rename_folder")	
	}
	if ok, err := es.EmailRepo.CheckFolder(email, newFolderName); ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("folder_already_exists")	
	}
	return es.EmailRepo.RenameFolder(email, folderName, newFolderName)
}

func (es *EmailService) ChangeEmailFolder(id int, email string, folderName string) error {
	return es.EmailRepo.ChangeEmailFolder(id, email, folderName)
}


func (es *EmailService) CreateDraft(email models.Email) error {
	return es.EmailRepo.CreateDraft(email)
}

func (es *EmailService) UpdateDraft(email models.Draft) error {
	return es.EmailRepo.UpdateDraft(email)
}

func (es *EmailService) SendDraft(email models.Email) error {
	err := es.EmailRepo.DeleteEmails(email.Sender_email, []int{email.ID})
	if err != nil {
		return err
	}
	return es.EmailRepo.SaveEmail(email)
}