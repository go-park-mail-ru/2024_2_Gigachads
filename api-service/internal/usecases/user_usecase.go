package usecase

import (

	"context"
	"fmt"
	"net/http"
	models "mail/api-service/internal/models"
	"mail/api-service/pkg/utils"
	"os"	
)

type UserService struct {
	UserRepo    models.UserRepository
}

func NewUserService(urepo models.UserRepository) models.UserUseCase {
	return &UserService{
		UserRepo:    urepo,
	}
}

func (us *UserService) ChangeAvatar(fileContent []byte, email string) error {
	if err := os.MkdirAll("./avatars", os.ModePerm); err != nil {
		return err
	}

	var fileExtension string
	fileMIMEType := http.DetectContentType(fileContent)
	switch fileMIMEType {
	case "image/jpeg", "image/png":
		fileExtension = fileMIMEType[6:]
	default:
		return fmt.Errorf("do_not_support_mime_type_%s", fileMIMEType)
	}

	user, err := us.UserRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	var fileName string
	if fileName, err = utils.GenerateHash(); err != nil {
		return err
	}

	fileName += "." + fileExtension
	filePath := "./avatars/" + fileName

	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = outFile.Write(fileContent)
	if err != nil {
		return err
	}

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