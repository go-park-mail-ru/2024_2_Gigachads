package usecase

import (

	"context"
	"fmt"
	"io"
	models "mail/internal/models"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

type UserService struct {
	UserRepo    models.UserRepository
	SessionRepo models.SessionRepository
	CsrfRepo    models.CsrfRepository
}

func NewUserService(urepo models.UserRepository, srepo models.SessionRepository, crepo models.CsrfRepository) models.UserUseCase {
	return &UserService{
		UserRepo:    urepo,
		SessionRepo: srepo,
		CsrfRepo:    crepo,
	}
}

func (us *UserService) Signup(ctx context.Context, signup *models.User) (*models.User, *models.Session, *models.Csrf, error) {
	taken, err := us.UserRepo.IsExist(signup.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	if taken {
		return nil, nil, nil, fmt.Errorf("login_taken")
	}

	user, err := us.UserRepo.CreateUser(signup)
	if err != nil {
		return nil, nil, nil, err
	}

	session, err := us.SessionRepo.CreateSession(ctx, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	csrf, err := us.CsrfRepo.CreateCsrf(ctx, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	return user, session, csrf, nil
}

func (us *UserService) Login(ctx context.Context, login *models.User) (*models.User, *models.Session, *models.Csrf, error) {
	taken, err := us.UserRepo.IsExist(login.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	if !taken {
		return nil, nil, nil, fmt.Errorf("user_does_not_exist")
	}

	user, err := us.UserRepo.CheckUser(login)
	if err != nil {
		return nil, nil, nil, err
	}
	session, err := us.SessionRepo.CreateSession(ctx, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	csrf, err := us.CsrfRepo.CreateCsrf(ctx, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	return user, session, csrf, nil
}

func (us *UserService) Logout(ctx context.Context, id string) error {
	err := us.SessionRepo.DeleteSession(ctx, id)
	if err != nil {
		return err
	}
	err = us.CsrfRepo.DeleteCsrf(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) CheckAuth(ctx context.Context, id string) (string, error) {
	session, err := us.SessionRepo.GetSession(ctx, id)
	if err != nil {
		return "", err
	}
	return session, nil
}

func (us *UserService) CheckCsrf(ctx context.Context, session string, csrf string) error {
	email1, err := us.CsrfRepo.GetCsrf(ctx, csrf)
	if err != nil {
		return err
	}
	email2, err := us.SessionRepo.GetSession(ctx, session)
	if err != nil {
		return err
	}
	if email1 != email2 {
		return fmt.Errorf("invalid_csrf")
	}
	return nil
}

func (us *UserService) ChangeAvatar(file multipart.File, header multipart.FileHeader, email string) error {
	if err := os.MkdirAll("./avatars", os.ModePerm); err != nil {
		return err
	}

	if header.Size > (5 * 1024 * 1024) {
		return fmt.Errorf("too_big_file")
	}

	ext := filepath.Ext(header.Filename)
	user, err := us.UserRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	fileName := strconv.Itoa(user.ID) + ext
	filePath := "./avatars/" + fileName
	_, err = os.Stat(filePath)
	// if os.IsNotExist(err) {
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, file)
	if err != nil {
		return err
	}
	// } else {
	// 	outFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer outFile.Close()
	// 	_, err = io.Copy(outFile, file)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	user.AvatarURL = fileName
	err = us.UserRepo.UpdateInfo(user)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetAvatar(email string) ([]byte, string, error) {
	user, err := us.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, "", err
	}
	
	filePath := "./avatars/" + user.AvatarURL
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}
	return data, user.AvatarURL, nil
	
}

func (us *UserService) ChangePassword(email string, password string) error {
	user, err := us.UserRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	user.Password = password
	err = us.UserRepo.UpdateInfo(user)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) ChangeName(email string, name string) error {
	user, err := us.UserRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	user.Name = name
	err = us.UserRepo.UpdateInfo(user)
	if err != nil {
		return err
	}
	return nil
}
