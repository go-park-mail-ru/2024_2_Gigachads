package middleware

import (
	"context"
	"mail/internal/models"
	"net/http"
)

type AuthMiddleware struct {
	UserUseCase models.UserUseCase
}

func NewAuthMW(uu models.UserUseCase) *AuthMiddleware {
	return &AuthMiddleware{UserUseCase: uu}
}

type contextKey string

const Key = contextKey("email")

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("email")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		email, err := m.UserUseCase.CheckAuth(r.Context(), cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), Key, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
