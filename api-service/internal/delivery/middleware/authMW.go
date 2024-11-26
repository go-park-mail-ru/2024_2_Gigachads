package middleware

import (
	"context"
	"mail/models"
	"net/http"
)

type AuthMiddleware struct {
	AuthUseCase models.AuthUseCase
}

func NewAuthMW(au models.AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{AuthUseCase: au}
}

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
		sessionID := ""
		if cookie != nil {
			sessionID = cookie.Value
		}

		cookie, err = r.Cookie("csrf")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		csrf := ""
		if cookie != nil {
			csrf = cookie.Value
		}

		email, err := m.AuthUseCase.CheckAuth(r.Context(), sessionID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		err = m.AuthUseCase.CheckCsrf(r.Context(), sessionID, csrf)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "email", email)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
