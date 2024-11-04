package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"mail/pkg/utils"
	"mail/pkg/logger"
	models "mail/internal/models"
)

type UserRepositoryService struct {
	repo *sql.DB
	logger logger.Logable
}

func NewUserRepositoryService(db *sql.DB, l logger.Logable) models.UserRepository {
	return &UserRepositoryService{repo: db, logger: l}
}

func (ur *UserRepositoryService) GetByEmail(email string) (bool, error) {

	email = utils.Sanitize(email)

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
		ur.logger.Error(err.Error())
		return false, err
	}
	return false, nil
}

func (ur *UserRepositoryService) CreateUser(signup *models.User) (*models.User, error) {

	signup.Name = utils.Sanitize(signup.Name)
	signup.Email = utils.Sanitize(signup.Email)
	signup.Password = utils.Sanitize(signup.Password)

	row := ur.repo.QueryRow(
		`INSERT INTO "profile" (username, email, password) VALUES ($1, $2, $3) RETURNING email`,
		signup.Name, signup.Email, signup.Password)
	user := models.User{}
	err := row.Scan(&user.Email)
	if err != nil {
		ur.logger.Error(err.Error())
		return nil, err
	}

	user.Name = utils.Sanitize(user.Name)
	user.Email = utils.Sanitize(user.Email)
	user.Password = utils.Sanitize(user.Password)

	return &user, nil
}

func (ur *UserRepositoryService) CheckUser(login *models.User) (*models.User, error) {

	login.Name = utils.Sanitize(login.Name)
	login.Email = utils.Sanitize(login.Email)
	login.Password = utils.Sanitize(login.Password)

	row := ur.repo.QueryRow(
		`SELECT email, password FROM "user" WHERE email = $1`, login.Email)
	user := models.User{}
	err := row.Scan(&user.Email, &user.Password)
	if err != nil {
		ur.logger.Error(err.Error())
		return nil, err
	}

	if login.Password != user.Password {
		return nil, fmt.Errorf("invalid_password")
	}

	user.Name = utils.Sanitize(user.Name)
	user.Email = utils.Sanitize(user.Email)
	user.Password = utils.Sanitize(user.Password)

	return &user, nil
}
