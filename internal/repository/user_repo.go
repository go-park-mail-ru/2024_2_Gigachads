package repository

import (
	"fmt"
	models "mail/internal/models"
)

type UserRepository struct {
	repo map[string]models.User
}

func NewUserRepository() *UserRepository {
	repo := make(map[string]models.User)
	return &UserRepository{repo: repo}
}

func (ur *UserRepository) CreateUser(signup *models.Signup) (*models.User, error) {
	_, ok := ur.repo[signup.Email]
	if ok {
		return &models.User{}, fmt.Errorf("login_taken")
	}
	user := models.User{Name: signup.Name, Email: signup.Email, Password: signup.Password}
	ur.repo[user.Email] = user
	return &user, nil
}

func (ur *UserRepository) GetUser(login *models.Login) (*models.User, error) {
	user, ok := ur.repo[login.Email]
	if ok {
		if user.Password != login.Password {
			return &models.User{}, fmt.Errorf("invalid_password")
		}
		return &user, nil
	} else {
		return &models.User{}, fmt.Errorf("user_does_not_exist")
	}
}
