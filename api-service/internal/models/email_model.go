package models

import (
	"context"
	"time"
	"encoding/json"
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
	Attachments  []string  `json:"attachments"`
	Files        []File    `json:"filenames"`
}

// type Draft struct {
// 	ID          int    `json:"id"`
// 	Recipient   string `json:"recipient"`
// 	Title       string `json:"title"`
// 	Description string `json:"description"`
// 	ParentID    int    `json:"parentID"`
// }

type Folder struct {
	Name string `json:"name"`
}

type RenameFolder struct {
	Name    string `json:"name"`
	NewName string `json:"new_name"`
}

type Timestamp struct {
	LastModified time.Time `json:"date"`
}

func (u Timestamp) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u Timestamp) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

type FilePath struct {
	Path string `json:"path"`
}

type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type SmtpPop3Usecase interface {
	SendEmail(ctx context.Context, from string, to []string, subject string, body string) error
	ForwardEmail(ctx context.Context, from string, to []string, originalEmail Email) error
	ReplyEmail(ctx context.Context, from string, to string, originalEmail Email, replyText string) error
}

type EmailUseCase interface {
	Inbox(id string) ([]Email, error)
	SendEmail(ctx context.Context, from string, to []string, subject string, body string) error
	ForwardEmail(ctx context.Context, from string, to []string, originalEmail Email) error
	ReplyEmail(ctx context.Context, from string, to string, originalEmail Email, replyText string) error
	GetEmailByID(id int) (Email, error)
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
	UpdateDraft(email Email) error
	SendDraft(email Email) error
	InboxStatus(ctx context.Context, email string, lastModified time.Time) ([]Email, error)
	UploadAttach(ctx context.Context, fileContent []byte, filename string) (string, error)
	GetAttach(ctx context.Context, path string) ([]byte, error)
	DeleteAttach(ctx context.Context, path string) error
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
	GetNewEmails(email string, LastModified time.Time) ([]Email, error)
	CreateFolder(email string, folderName string) error
	DeleteFolder(email string, folderName string) error
	RenameFolder(email string, folderName string, newFolderName string) error
	ChangeEmailFolder(id int, email string, folderName string) error
	CreateDraft(email Email) error
	UpdateDraft(email Email) error
	CheckFolder(email string, folderName string) (bool, error)
	GetMessageFolder(msgID int) (string, error)
	GetTimestamp(ctx context.Context, email string) (time.Time, error)
	SetTimestamp(ctx context.Context, email string) error
	DeleteAttach(ctx context.Context, path string) error
	GetAttach(ctx context.Context, path string) ([]byte, error)
	UploadAttach(ctx context.Context, fileContent []byte, filename string) (string, error)
	ConnectAttachToMessage(messageID int, path string) error
}
