package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"mail/internal/models"
	"net/http"
)

type AuthMiddleware struct {
	UserUseCase models.UserUseCase
}

func NewAuthMW(uu models.UserUseCase) *AuthMiddleware {
	return &AuthMiddleware{UserUseCase: uu}
}

func (mw *AuthMiddleware) ConfigureAuthMiddleware(privateMux *mux.Router) {
	privateMux.Use(mw.Handler)
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
