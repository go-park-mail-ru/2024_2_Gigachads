package repository

import (
	"mail/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryService_IsExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepositoryService(db)

	tests := []struct {
		name    string
		email   string
		mock    func()
		want    bool
		wantErr bool
	}{
		{
			name:  "пользователь существует",
			email: "test@example.com",
			mock: func() {
				rows := sqlmock.NewRows([]string{"email"}).AddRow("test@example.com")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email FROM "profile" WHERE email = $1`)).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.IsExist(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepositoryService_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepositoryService(db)

	tests := []struct {
		name    string
		user    *models.User
		mock    func()
		want    *models.User
		wantErr bool
	}{
		{
			name: "успешное создание пользователя",
			user: &models.User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"email"}).AddRow("test@example.com")
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "profile"`)).
					WithArgs("Test User", "test@example.com", "password123").
					WillReturnRows(rows)
			},
			want: &models.User{
				Email: "test@example.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.CreateUser(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUserRepositoryService_CheckUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepositoryService(db)

	tests := []struct {
		name    string
		login   *models.User
		mock    func()
		want    *models.User
		wantErr bool
	}{
		{
			name: "успешная проверка пользователя",
			login: &models.User{
				Email:    "test@example.com",
				Password: "password123",
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"email", "username", "password"}).
					AddRow("test@example.com", "Test User", "password123")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, username, password FROM "profile"`)).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			want: &models.User{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.CheckUser(tt.login)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUserRepositoryService_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepositoryService(db)

	tests := []struct {
		name    string
		email   string
		mock    func()
		want    *models.User
		wantErr bool
	}{
		{
			name:  "успешное получение пользователя",
			email: "test@example.com",
			mock: func() {
				rows := sqlmock.NewRows([]string{"email", "id", "password", "username", "avatar_url"}).
					AddRow("test@example.com", 1, "password123", "Test User", "avatar.jpg")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, id, password, username, avatar_url FROM "profile"`)).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			want: &models.User{
				Email:     "test@example.com",
				ID:        1,
				Password:  "password123",
				Name:      "Test User",
				AvatarURL: "avatar.jpg",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.GetUserByEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUserRepositoryService_UpdateInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepositoryService(db)

	tests := []struct {
		name    string
		user    *models.User
		mock    func()
		wantErr bool
	}{
		{
			name: "успешное обновление информации",
			user: &models.User{
				Email:     "test@example.com",
				Name:      "Updated Name",
				Password:  "newpassword",
				AvatarURL: "new_avatar.jpg",
			},
			mock: func() {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "profile"`)).
					WithArgs("Updated Name", "new_avatar.jpg", "newpassword", "test@example.com").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.UpdateInfo(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
