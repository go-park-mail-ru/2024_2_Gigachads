package repository

import (
	"database/sql"
	"mail/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryService_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	repo := NewUserRepositoryService(db)

	t.Run("успешное создание пользователя", func(t *testing.T) {
		user := &models.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password",
		}

		rows := sqlmock.NewRows([]string{"email"}).AddRow(user.Email)
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "profile" (username, email, password) VALUES ($1, $2, $3) RETURNING email`)).
			WithArgs(user.Name, user.Email, user.Password).
			WillReturnRows(rows)

		createdUser, err := repo.CreateUser(user)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, createdUser.Email)
	})

	t.Run("ошибка создания пользователя", func(t *testing.T) {
		user := &models.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "profile"`)).
			WithArgs(user.Name, user.Email, user.Password).
			WillReturnError(sql.ErrConnDone)

		createdUser, err := repo.CreateUser(user)
		assert.Error(t, err)
		assert.Nil(t, createdUser)
	})
}

func TestUserRepositoryService_CheckUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	repo := NewUserRepositoryService(db)

	t.Run("успешная проверка пользователя", func(t *testing.T) {
		user := &models.User{
			Email:    "test@example.com",
			Password: "password",
		}

		rows := sqlmock.NewRows([]string{"email", "username", "password"}).
			AddRow(user.Email, "Test User", user.Password)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, username, password FROM "profile" WHERE email = $1`)).
			WithArgs(user.Email).
			WillReturnRows(rows)

		checkedUser, err := repo.CheckUser(user)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, checkedUser.Email)
		assert.Equal(t, user.Password, checkedUser.Password)
	})

	t.Run("неверный пароль", func(t *testing.T) {
		user := &models.User{
			Email:    "test@example.com",
			Password: "wrong_password",
		}

		rows := sqlmock.NewRows([]string{"email", "username", "password"}).
			AddRow(user.Email, "Test User", "correct_password")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, username, password FROM "profile" WHERE email = $1`)).
			WithArgs(user.Email).
			WillReturnRows(rows)

		checkedUser, err := repo.CheckUser(user)
		assert.Error(t, err)
		assert.Equal(t, "invalid_password", err.Error())
		assert.Nil(t, checkedUser)
	})

	t.Run("пользователь не найден", func(t *testing.T) {
		user := &models.User{
			Email:    "nonexistent@example.com",
			Password: "password",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, username, password FROM "profile" WHERE email = $1`)).
			WithArgs(user.Email).
			WillReturnError(sql.ErrNoRows)

		checkedUser, err := repo.CheckUser(user)
		assert.Error(t, err)
		assert.Nil(t, checkedUser)
	})
}

func TestUserRepositoryService_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	repo := NewUserRepositoryService(db)

	t.Run("успешное получение пользователя", func(t *testing.T) {
		email := "test@example.com"
		expectedUser := &models.User{
			Email:     email,
			ID:        1,
			Password:  "password",
			Name:      "Test User",
			AvatarURL: "avatar.jpg",
		}

		rows := sqlmock.NewRows([]string{"email", "id", "password", "username", "avatar_url"}).
			AddRow(expectedUser.Email, expectedUser.ID, expectedUser.Password, expectedUser.Name, expectedUser.AvatarURL)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, id, password, username, avatar_url FROM "profile" WHERE email = $1`)).
			WithArgs(email).
			WillReturnRows(rows)

		user, err := repo.GetUserByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("пользователь не найден", func(t *testing.T) {
		email := "nonexistent@example.com"

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, id, password, username, avatar_url FROM "profile" WHERE email = $1`)).
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByEmail(email)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepositoryService_UpdateInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	repo := NewUserRepositoryService(db)

	t.Run("успешное обновление информации", func(t *testing.T) {
		user := &models.User{
			Email:     "test@example.com",
			Name:      "New Name",
			AvatarURL: "new_avatar.jpg",
			Password:  "new_password",
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "profile" SET username = $1, avatar_url = $2, password = $3 WHERE email = $4`)).
			WithArgs(user.Name, user.AvatarURL, user.Password, user.Email).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateInfo(user)
		assert.NoError(t, err)
	})

	t.Run("ошибка обновления", func(t *testing.T) {
		user := &models.User{
			Email: "test@example.com",
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "profile"`)).
			WithArgs(user.Name, user.AvatarURL, user.Password, user.Email).
			WillReturnError(sql.ErrNoRows)

		err := repo.UpdateInfo(user)
		assert.Error(t, err)
	})
}
