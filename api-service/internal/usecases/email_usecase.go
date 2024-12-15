package usecase

import (
	"context"
	"fmt"
	"time"
	"errors"
	"database/sql"
	"path/filepath"
	"google.golang.org/protobuf/types/known/timestamppb"
	"mail/api-service/internal/models"
	proto "mail/gen/go/smtp"
)

type EmailService struct {
	EmailRepo models.EmailRepository
	EmailMS   proto.SmtpPop3ServiceClient
}

func NewEmailService(
	erepo models.EmailRepository,
	client proto.SmtpPop3ServiceClient,
) *EmailService {
	return &EmailService{
		EmailRepo: erepo,
		EmailMS:   client,
	}
}

func (es *EmailService) Inbox(email string) ([]models.Email, error) {
	// return es.EmailRepo.Inbox(email)
	return es.EmailRepo.GetFolderEmails(email, "Входящие")
}

func (es *EmailService) GetEmailByID(id int) (models.Email, error) {
	email, err :=  es.EmailRepo.GetEmailByID(id)
		if err != nil {
		return models.Email{}, err
	}

	for _, path := range email.Attachments {
		var filebody models.File
		file := filepath.Base(path)
		filebody.Path = path
		filebody.Name = file
		email.Files = append(email.Files, filebody)
	}
	return email, nil
}

func (es *EmailService) GetSentEmails(email string) ([]models.Email, error) {
	//return es.EmailRepo.GetSentEmails(senderEmail)
	return es.EmailRepo.GetFolderEmails(email, "Отправленные")
}

func (s *EmailService) SaveEmail(email models.Email) error {
	err := s.EmailRepo.SaveEmail(email)
	if err == sql.ErrNoRows {
		return errors.New("email_not_found")
	}
	if err != nil {
		return errors.New("failed_to_save_email")
	}
	return nil
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
	if folderName == "Отправленные" || folderName == "Черновики" {
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

func (es *EmailService) UpdateDraft(email models.Email) error {
	return es.EmailRepo.UpdateDraft(email)
}

func (es *EmailService) SendDraft(email models.Email) error {
	m, err := es.EmailRepo.GetEmailByID(email.ID)
	if err != nil {
		return err
	}
	err = es.EmailRepo.DeleteEmails(email.Sender_email, []int{email.ID})
	if err != nil {
		return err
	}
	return es.EmailRepo.SaveEmail(m)
}

func (es *EmailService) SendEmail(ctx context.Context, from string, to []string, subject string, body string) error {
	for i := range to {
		req := &proto.SendEmailRequest{From: from, To: to[i], Subject: subject, Body: body}
		_, err := es.EmailMS.SendEmail(ctx, req)
		if err != nil {
			return err
		}
		es.EmailRepo.SetTimestamp(ctx, to[i])
	}
	return nil
}

func (es *EmailService) ForwardEmail(ctx context.Context, from string, to []string, originalEmail models.Email) error {
	for i := range to {
		req := &proto.ForwardEmailRequest{SendingDate: timestamppb.New(originalEmail.Sending_date), Sender: originalEmail.Sender_email, From: from, To: to[i], Title: originalEmail.Title, Description: originalEmail.Description}
		_, err := es.EmailMS.ForwardEmail(ctx, req)
		if err != nil {
			return err
		}
		es.EmailRepo.SetTimestamp(ctx, to[i])
	}
	return nil
}

func (es *EmailService) ReplyEmail(ctx context.Context, from string, to string, originalEmail models.Email, replyText string) error {
	req := &proto.ReplyEmailRequest{ReplyText: replyText, SendingDate: timestamppb.New(originalEmail.Sending_date), Sender: originalEmail.Sender_email, From: from, To: to, Title: originalEmail.Title, Description: originalEmail.Description}
	_, err := es.EmailMS.ReplyEmail(ctx, req)
	es.EmailRepo.SetTimestamp(ctx, to)
	return err
}

func (es *EmailService) InboxStatus(ctx context.Context, email string, frontLastModified time.Time) ([]models.Email, error) {
	lastModified, err := es.EmailRepo.GetTimestamp(ctx, email)
	if err != nil {
		return nil, err
	}
	if (frontLastModified.Before(lastModified)){
		return es.EmailRepo.GetNewEmails(email, frontLastModified)
	} else {
		return nil, fmt.Errorf("not_modified")
	}

	return nil, nil
}

func (es *EmailService) UploadAttach(fileContent []byte, filename string) (string, error) {
	return es.EmailRepo.UploadAttach(fileContent, filename)
}

func (es *EmailService) GetAttach(path string) ([]byte, error) {
	return es.EmailRepo.GetAttach(path)
}

func (es *EmailService) DeleteAttach(path string) error {
	return es.EmailRepo.DeleteAttach(path)
}