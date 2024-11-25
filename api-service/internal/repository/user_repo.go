package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"mail/models"
)

type UserRepositoryService struct {
	repo *sql.DB
}

func NewUserRepositoryService(db *sql.DB) models.UserRepository {
	return &UserRepositoryService{repo: db}
}

func (ur *UserRepositoryService) IsExist(email string) (bool, error) {
	row := ur.repo.QueryRow(
		`SELECT email FROM "profile" WHERE email = $1`, email)
	user := models.User{}
	err := row.Scan(&user.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
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
		`SELECT email, username, password FROM "profile" WHERE email = $1`, login.Email)
	user := models.User{}
	err := row.Scan(&user.Email, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}

	if login.Password != user.Password {
		return nil, fmt.Errorf("invalid_password")
	}

	return &user, nil
}

func (ur *UserRepositoryService) GetUserByEmail(email string) (*models.User, error) {
	row := ur.repo.QueryRow(
		`SELECT email, id, password, username, avatar_url FROM "profile" WHERE email = $1`, email)
	user := models.User{}
	err := row.Scan(&user.Email, &user.ID, &user.Password, &user.Name, &user.AvatarURL)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepositoryService) UpdateInfo(user *models.User) error {
	query :=
		`UPDATE "profile"
		SET username = $1, avatar_url = $2, password = $3
		WHERE email = $4`
	_, err := ur.repo.Exec(query, user.Name, user.AvatarURL, user.Password, user.Email)
	if err != nil {
		return err
	}
	return nil
}
