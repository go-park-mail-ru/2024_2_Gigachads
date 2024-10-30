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

		cookie, err := r.Cookie("session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		session, err := m.UserUseCase.CheckAuth(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), Key, session.UserLogin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
