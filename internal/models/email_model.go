package models

import (
	"time"
)

type Email struct {
	ID           int       `json:"id"`
	ParentID     int       `json:"parentID"`
	Sender_email string    `json:"sender"`
	Recipient    string    `json:"recipient"`
	Title        string    `json:"title"`
	IsRead       bool      `json:"isRead"`
	Sending_date time.Time `json:"date"`
	Description  string    `json:"description"`
}

type Draft struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
}

type Folder struct {
	Name 		 string    `json:"name"`
}

type RenameFolder struct {
	Name 		 string    `json:"name"`
	NewName 	 string    `json:"new_name"`
}

type EmailUseCase interface {
	Inbox(id string) ([]Email, error)
	SendEmail(from string, to []string, subject string, body string) error
	ForwardEmail(from string, to []string, originalEmail Email) error
	ReplyEmail(from string, to string, originalEmail Email, replyText string) error
	GetEmailByID(id int) (Email, error)
	FetchEmailsViaPOP3() error
	ChangeStatus(id int, status bool) error
	GetSentEmails(senderEmail string) ([]Email, error)
	SaveEmail(email Email) error
	DeleteEmails(userEmail string, messageIDs []int) error
	GetFolders(email string) ([]string, error)
	GetFolderEmails(email string, folderName string) ([]Email, error)
	CreateFolder(email string, folderName string) error
	DeleteFolder(email string, folderName string) error
	RenameFolder(email string, folderName string, newFolderName string) error
	ChangeEmailFolder(id int, email string, folderName string) error
	CreateDraft(email Email) error
	UpdateDraft(email Draft) error
	SendDraft(email Email) error
}

type EmailRepository interface {
	Inbox(id string) ([]Email, error)
	GetEmailByID(id int) (Email, error)
	SaveEmail(email Email) error
	ChangeStatus(id int, status bool) error
	GetSentEmails(senderEmail string) ([]Email, error)
	DeleteEmails(userEmail string, messageIDs []int) error
	GetFolders(email string) ([]string, error)
	GetFolderEmails(email string, folderName string) ([]Email, error)
	CreateFolder(email string, folderName string) error
	DeleteFolder(email string, folderName string) error
	RenameFolder(email string, folderName string, newFolderName string) error
	ChangeEmailFolder(id int, email string, folderName string) error
	CreateDraft(email Email) error
	UpdateDraft(email Draft) error
	CheckFolder(email string, folderName string) (bool, error)
}

type SMTPRepository interface {
	SendEmail(from string, to []string, subject string, body string) error
}

type POP3Repository interface {
	Connect() error
	FetchEmails(EmailRepository) error
	Quit() error
}
