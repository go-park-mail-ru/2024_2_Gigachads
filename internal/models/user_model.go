package models

import (
	"net/mail"
	"regexp"
)

type User struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	RePassword string `json:"repassword"`
}

type UserUseCase interface {
	Signup(user *User) (*User, *Session, error)
	Login(user *User) (*User, *Session, error)
	Logout(id string) error
	CheckAuth(sessionID string) (*Session, error)
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
