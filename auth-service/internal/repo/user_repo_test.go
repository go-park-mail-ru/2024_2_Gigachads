package repository

import (
	"database/sql"
	"mail/auth-service/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type MockLogger struct{}

func (m *MockLogger) Info(msg string, args ...any)  {}
func (m *MockLogger) Error(msg string, args ...any) {}
func (m *MockLogger) Fatal(msg string, args ...any) {}
func (m *MockLogger) Debug(msg string, args ...any) {}
func (m *MockLogger) Warn(msg string, args ...any)  {}

func TestIsExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	logger := &MockLogger{}
	repo := NewUserRepositoryService(db, logger)

	t.Run("пользователь существует", func(t *testing.T) {
		email := "test@test.com"
		rows := sqlmock.NewRows([]string{"email"}).AddRow(email)
		mock.ExpectQuery(`SELECT email FROM "profile"`).
			WithArgs(email).
			WillReturnRows(rows)

		exists, err := repo.IsExist(email)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("пользователь не существует", func(t *testing.T) {
		email := "nonexistent@test.com"
		mock.ExpectQuery(`SELECT email FROM "profile"`).
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		exists, err := repo.IsExist(email)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	logger := &MockLogger{}
	repo := NewUserRepositoryService(db, logger)

	t.Run("успешное создание", func(t *testing.T) {
		user := &models.User{
			Email:    "test@test.com",
			Name:     "Test User",
			Password: "password",
		}

		rows := sqlmock.NewRows([]string{"email"}).AddRow(user.Email)
		mock.ExpectQuery(`INSERT INTO "profile"`).
			WithArgs(user.Name, user.Email, user.Password).
			WillReturnRows(rows)

		created, err := repo.CreateUser(user)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, created.Email)
	})

	t.Run("ошибка создания", func(t *testing.T) {
		user := &models.User{
			Email:    "test@test.com",
			Name:     "Test User",
			Password: "password",
		}

		mock.ExpectQuery(`INSERT INTO "profile"`).
			WithArgs(user.Name, user.Email, user.Password).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.CreateUser(user)
		assert.Error(t, err)
	})
}

func TestCheckUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	logger := &MockLogger{}
	repo := NewUserRepositoryService(db, logger)

	t.Run("успешная проверка", func(t *testing.T) {
		login := &models.User{
			Email:    "test@test.com",
			Password: "password",
		}

		rows := sqlmock.NewRows([]string{"email", "username", "password"}).
			AddRow(login.Email, "Test User", login.Password)
		mock.ExpectQuery(`SELECT email, username, password FROM "profile"`).
			WithArgs(login.Email).
			WillReturnRows(rows)

		user, err := repo.CheckUser(login)
		assert.NoError(t, err)
		assert.Equal(t, login.Email, user.Email)
		assert.Equal(t, login.Password, user.Password)
	})
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	logger := &MockLogger{}
	repo := NewUserRepositoryService(db, logger)

	t.Run("успешное получение", func(t *testing.T) {
		email := "test@test.com"
		rows := sqlmock.NewRows([]string{"email", "id", "password", "username", "avatar_url"}).
			AddRow(email, 1, "password", "Test User", "avatar.jpg")
		mock.ExpectQuery(`SELECT email, id, password, username, avatar_url FROM "profile"`).
			WithArgs(email).
			WillReturnRows(rows)

		user, err := repo.GetUserByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "avatar.jpg", user.AvatarURL)
	})
}

func TestUpdateInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	logger := &MockLogger{}
	repo := NewUserRepositoryService(db, logger)

	t.Run("успешное обновление", func(t *testing.T) {
		user := &models.User{
			Email:     "test@test.com",
			Name:      "New Name",
			AvatarURL: "new_avatar.jpg",
			Password:  "new_password",
		}

		mock.ExpectExec(`UPDATE "profile"`).
			WithArgs(user.Name, user.AvatarURL, user.Password, user.Email).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateInfo(user)
		assert.NoError(t, err)
	})

	t.Run("ошибка обновления", func(t *testing.T) {
		user := &models.User{
			Email:     "test@test.com",
			Name:      "New Name",
			AvatarURL: "new_avatar.jpg",
			Password:  "new_password",
		}

		mock.ExpectExec(`UPDATE "profile"`).
			WithArgs(user.Name, user.AvatarURL, user.Password, user.Email).
			WillReturnError(sql.ErrConnDone)

		err := repo.UpdateInfo(user)
		assert.Error(t, err)
	})
}
