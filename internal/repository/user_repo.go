package repository

import (
	"database/sql"
	"errors"
	"fmt"
	models "mail/internal/models"
)

type UserRepositoryService struct {
	repo *sql.DB
}

func NewUserRepositoryService(db *sql.DB) models.UserRepository {
	return &UserRepositoryService{repo: db}
}

func (ur *UserRepositoryService) GetByEmail(email string) (bool, error) {
	row := ur.repo.QueryRow(
		`SELECT email FROM "profile" WHERE email = $1`, email)
	user := models.User{}
	err := row.Scan(&user.Email)

	if !errors.Is(err, sql.ErrNoRows) {
		return true, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return false, nil
}

func (ur *UserRepositoryService) CreateUser(signup *models.User) (*models.User, error) {
	row := ur.repo.QueryRow(
		`INSERT INTO "profile" (username, email, password) VALUES ($1, $2, $3) RETURNING email`,
		signup.Name, signup.Email, signup.Password)
	user := models.User{}
	err := row.Scan(&user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepositoryService) CheckUser(login *models.User) (*models.User, error) {
	row := ur.repo.QueryRow(
		`SELECT email, password FROM "user" WHERE email = $1`, login.Email)
	user := models.User{}
	err := row.Scan(&user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	if login.Password != user.Password {
		return nil, fmt.Errorf("invalid_password")
	}

	return &user, nil
}
