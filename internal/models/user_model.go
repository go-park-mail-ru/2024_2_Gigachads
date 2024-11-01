package models

import (
	"net/mail"
	"regexp"
	"context"
)

type User struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	RePassword string `json:"repassword"`
}

type UserUseCase interface {
	Signup(ctx context.Context, user *User) (*User, *Session, error)
	Login(ctx context.Context, user *User) (*User, *Session, error)
	Logout(ctx context.Context, id string) error
	CheckAuth(ctx context.Context, sessionID string) (string, error)
}

type UserRepository interface {
	CreateUser(signup *User) (*User, error)
	CheckUser(login *User) (*User, error)
	GetByEmail(email string) (bool, error)
}

func EmailIsValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func InputIsValid(str string) bool {
	match, err := regexp.MatchString("^[a-zA-Z0-9_]+$", str)
	return match || err == nil
}
