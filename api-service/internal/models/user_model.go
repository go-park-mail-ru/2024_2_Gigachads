package models

import (
	"context"
	"net/mail"
	"regexp"
)

type User struct {
	ID         int
	Email      string `json:"email"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	RePassword string `json:"repassword"`
	AvatarURL  string `json:"avatarPath"`
}
type ChangeName struct {
	Name string `json:"name"`
}

type UserLogin struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarPath"`
}

type ChangePassword struct {
	Password    string `json:"password"`
	OldPassword string `json:"oldpassword"`
	RePassword  string `json:"repassword"`
}

type UserUseCase interface {
	ChangePassword(email string, password string) error
	ChangeName(email string, name string) error
	GetAvatar(email string) ([]byte, string, error)
	ChangeAvatar(fileContent []byte, email string) error
}

type UserRepository interface {
	CreateUser(signup *User) (*User, error)
	CheckUser(login *User) (*User, error)
	IsExist(email string) (bool, error)
	UpdateInfo(user *User) error
	GetUserByEmail(email string) (*User, error)
}

func EmailIsValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func InputIsValid(str string) bool {
	match, err := regexp.MatchString("^[a-zA-Z0-9_]+$", str)
	return match || err == nil
}

type AuthUseCase interface {
	Signup(ctx context.Context, signup *User) (string, string, error)
	Login(ctx context.Context, login *User) (string, string, string, string, error)
	Logout(ctx context.Context, email string) error
	CheckAuth(ctx context.Context, id string) (string, error)
	CheckCsrf(ctx context.Context, email string, csrf string) error
}