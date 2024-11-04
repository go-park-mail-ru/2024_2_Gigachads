package usecase

import (
	"context"
	"fmt"
	models "mail/internal/models"
	"path/filepath"
	"mime/multipart"
	"strconv"
	"os"
	"io"
	"bytes"
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
		CsrfRepo:	 crepo,
	}
}

<<<<<<< HEAD
func (us *UserService) Signup(ctx context.Context, signup *models.User) (*models.User, *models.Session, *models.Csrf, error) {
	taken, err := us.UserRepo.GetByEmail(signup.Email)
=======
func (us *UserService) Signup(ctx context.Context, signup *models.User) (*models.User, *models.Session, error) {
	taken, err := us.UserRepo.IsExist(signup.Email)
>>>>>>> cd849e7 (settings)
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
	taken, err := us.UserRepo.GetByEmail(login.Email)
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

<<<<<<< HEAD
func (us *UserService) CheckCsrf(ctx context.Context, id string) (string, error) {
	csrf, err := us.CsrfRepo.GetCsrf(ctx, id)
	if err != nil {
		return "", err
	}
	return csrf, nil
}
=======
func (us *UserService) ChangeAvatar(file multipart.File, header multipart.FileHeader, email string) (error) {
	if header.Size > (5 * 1024 * 1024) {
		return fmt.Errorf("too_big_file")
	}
	
	ext := filepath.Ext(header.Filename)
	user, err := us.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil
	}
	fileName := strconv.Itoa(user.ID) + "." + ext
	filePath := "./avatars/" + fileName
	
	_, err = os.Stat(filePath);
	if os.IsNotExist(err) {
		outFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, file)
		if err != nil {
			return err
		}
	} else {
		outFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, file)
		if err != nil {
			return err
		}
	}
	
	user.AvatarURL = fileName
	err = us.UserRepo.UpdateInfo(user)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetAvatar(email string) (*bytes.Buffer, error) {
	user, err := us.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, nil
	}
	
	filePath := "./avatars/" + user.AvatarURL
	
	file, err := os.Open(filePath)
	 if err != nil {
	  return nil, err
	 }
	 defer file.Close()
	
	 var buf bytes.Buffer
	 writer := multipart.NewWriter(&buf)
	
	 part, err := writer.CreateFormFile("file", filePath)
	 if err != nil {
	  return nil, err
	 }
	
	 if _, err := io.Copy(part, file); err != nil {
	  return nil, err
	 }

	 writer.Close()
	 return &buf, nil
}

func (us *UserService) ChangePassword(email string, password string) (error) {
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

func (us *UserService) ChangeName(email string, name string) (error) {
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
	


>>>>>>> cd849e7 (settings)
