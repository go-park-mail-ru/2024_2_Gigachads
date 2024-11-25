package models

import (
	"context"
)

type AuthUseCase interface {
	Signup(ctx context.Context, signup *User) (string, string, error)
	Login(ctx context.Context, login *User) (string, string, string, string, error)
	Logout(ctx context.Context, email string) error
	CheckAuth(ctx context.Context, id string) (string, error)
	CheckCsrf(ctx context.Context, email string, csrf string) error
}
