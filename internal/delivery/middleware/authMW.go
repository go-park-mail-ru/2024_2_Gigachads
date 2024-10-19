package middleware

import (
	"context"
	usecases "mail/internal/usecases"
	"net/http"
)

type AuthMiddleware struct {
	UserUseCase usecases.UserUseCase
}

func NewAuthMW(uu usecases.UserUseCase) *AuthMiddleware {
	return &AuthMiddleware{UserUseCase: uu}
}

type contextKey string

const Key = contextKey("session")

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		session, err := m.UserUseCase.CheckAuth(cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), Key, session.UserLogin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
