//go:generate mockgen -source=../../models/interfaces.go -destination=mocks/mock_interfaces.go -package=mocks

package usecase

import (
	"fmt"
	"mail/api-service/internal/models"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

func (m *MockUserRepository) CheckUser(user *models.User) (*models.User, error) {
	ret := m.ctrl.Call(m, "CheckUser", user)
	resultUser, _ := ret[0].(*models.User)
	var err error
	if ret[1] != nil {
		err = ret[1].(error)
	}
	return resultUser, err
}

func (m *MockUserRepositoryMockRecorder) CheckUser(user interface{}) *gomock.Call {
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "CheckUser", reflect.TypeOf((*MockUserRepository)(nil).CheckUser), user)
}

func (m *MockUserRepository) CreateUser(user *models.User) (*models.User, error) {
	ret := m.ctrl.Call(m, "CreateUser", user)
	resultUser, _ := ret[0].(*models.User)
	var err error
	if ret[1] != nil {
		err = ret[1].(error)
	}
	return resultUser, err
}

func (m *MockUserRepositoryMockRecorder) CreateUser(user interface{}) *gomock.Call {
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), user)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	user, _ := ret[0].(*models.User)
	var err error
	if ret[1] != nil {
		err = ret[1].(error)
	}
	return user, err
}

func (m *MockUserRepositoryMockRecorder) GetUserByEmail(email interface{}) *gomock.Call {
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "GetUserByEmail", reflect.TypeOf((*MockUserRepository)(nil).GetUserByEmail), email)
}

func (m *MockUserRepository) UpdateInfo(user *models.User) error {
	ret := m.ctrl.Call(m, "UpdateInfo", user)
	var err error
	if ret[0] != nil {
		err = ret[0].(error)
	}
	return err
}

func (m *MockUserRepositoryMockRecorder) UpdateInfo(user interface{}) *gomock.Call {
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "UpdateInfo", reflect.TypeOf((*MockUserRepository)(nil).UpdateInfo), user)
}

func (m *MockUserRepository) IsExist(email string) (bool, error) {
	ret := m.ctrl.Call(m, "IsExist", email)
	var err error
	if ret[1] != nil {
		err = ret[1].(error)
	}
	return ret[0].(bool), err
}

func (m *MockUserRepositoryMockRecorder) IsExist(email interface{}) *gomock.Call {
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "IsExist", reflect.TypeOf((*MockUserRepository)(nil).IsExist), email)
}

func TestUserService_ChangeAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	testEmail := "test@example.com"
	jpegContent := []byte{0xFF, 0xD8, 0xFF}
	invalidContent := []byte{0x00, 0x01, 0x02}

	tests := []struct {
		name        string
		email       string
		fileContent []byte
		mockFunc    func()
		wantErr     bool
	}{
		{
			name:        "Success JPEG avatar",
			email:       testEmail,
			fileContent: jpegContent,
			mockFunc: func() {
				mockRepo.EXPECT().
					GetUserByEmail(testEmail).
					Return(&models.User{Email: testEmail}, nil)
				mockRepo.EXPECT().
					UpdateInfo(gomock.Any()).
					AnyTimes().
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "Invalid file format",
			email:       testEmail,
			fileContent: invalidContent,
			mockFunc:    func() {},
			wantErr:     true,
		},
		{
			name:        "User not found",
			email:       testEmail,
			fileContent: jpegContent,
			mockFunc: func() {
				mockRepo.EXPECT().
					GetUserByEmail(testEmail).
					Return(nil, fmt.Errorf("user not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := service.ChangeAvatar(tt.fileContent, tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				files, _ := os.ReadDir("./avatars")
				assert.Greater(t, len(files), 0)
			}
		})
	}

	os.RemoveAll("./avatars")
}

func TestUserService_GetAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	testEmail := "test@example.com"
	testAvatar := "test.jpg"
	testContent := []byte("test content")

	os.MkdirAll("./avatars", os.ModePerm)
	os.WriteFile("./avatars/"+testAvatar, testContent, os.ModePerm)
	defer os.RemoveAll("./avatars")

	tests := []struct {
		name     string
		email    string
		mockFunc func()
		wantData []byte
		wantName string
		wantErr  bool
	}{
		{
			name:  "Success get avatar",
			email: testEmail,
			mockFunc: func() {
				mockRepo.EXPECT().
					GetUserByEmail(testEmail).
					Return(&models.User{Email: testEmail, AvatarURL: testAvatar}, nil)
			},
			wantData: testContent,
			wantName: testAvatar,
			wantErr:  false,
		},
		{
			name:  "User not found",
			email: testEmail,
			mockFunc: func() {
				mockRepo.EXPECT().
					GetUserByEmail(testEmail).
					Return(nil, fmt.Errorf("user not found"))
			},
			wantData: nil,
			wantName: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			data, name, err := service.GetAvatar(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantData, data)
				assert.Equal(t, tt.wantName, name)
			}
		})
	}
}

func TestUserService_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	testEmail := "test@example.com"
	testPassword := "password123"

	tests := []struct {
		name     string
		email    string
		password string
		mockFunc func()
		wantErr  bool
	}{
		{
			name:     "Success change password",
			email:    testEmail,
			password: testPassword,
			mockFunc: func() {
				mockRepo.EXPECT().
					GetUserByEmail(testEmail).
					Return(&models.User{Email: testEmail, Password: testPassword}, nil)
				mockRepo.EXPECT().
					UpdateInfo(gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Wrong password",
			email:    testEmail,
			password: "wrongpass",
			mockFunc: func() {
				mockRepo.EXPECT().
					GetUserByEmail(testEmail).
					Return(&models.User{Email: testEmail, Password: testPassword}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := service.ChangePassword(tt.email, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_ChangeName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	testEmail := "test@example.com"
	testName := "New Name"

	tests := []struct {
		name     string
		email    string
		newName  string
		mockFunc func()
		wantErr  bool
	}{
		{
			name:    "Success change name",
			email:   testEmail,
			newName: testName,
			mockFunc: func() {
				mockRepo.EXPECT().
					GetUserByEmail(testEmail).
					Return(&models.User{Email: testEmail}, nil)
				mockRepo.EXPECT().
					UpdateInfo(gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "User not found",
			email:   testEmail,
			newName: testName,
			mockFunc: func() {
				mockRepo.EXPECT().
					GetUserByEmail(testEmail).
					Return(nil, fmt.Errorf("user not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := service.ChangeName(tt.email, tt.newName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
