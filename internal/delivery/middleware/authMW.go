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

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		
		cookie, err := r.Cookie("email")
		sessionID := ""
		if cookie != nil {
			sessionID = cookie.Value
		}
		
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		
		email, err := m.UserUseCase.CheckAuth(r.Context(), sessionID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		
		ctx := context.WithValue(r.Context(), "email", email)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
