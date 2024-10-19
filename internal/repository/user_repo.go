package repository

import (
	"fmt"
	models "mail/internal/models"
)

type UserRepositoryService struct {
	repo map[string]models.User
}

func NewUserRepositoryService() models.UserRepository {
	repo := make(map[string]models.User)
	return &UserRepositoryService{repo: repo}
}

func (ur *UserRepositoryService) CreateUser(signup *models.User) (*models.User, error) {
	_, ok := ur.repo[signup.Email]
	if ok {
		return &models.User{}, fmt.Errorf("login_taken")
	}
	user := models.User{Name: signup.Name, Email: signup.Email, Password: signup.Password}
	ur.repo[user.Email] = user
	return &user, nil
}

func (ur *UserRepositoryService) CheckUser(login *models.User) (*models.User, error) {
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
