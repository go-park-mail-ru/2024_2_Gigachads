package auth

import (
	"mail/models"

	"github.com/gorilla/mux"
)

type AuthRouter struct {
	AuthUseCase models.AuthUseCase
}

func NewAuthRouter(au models.AuthUseCase) *AuthRouter {
	return &AuthRouter{AuthUseCase: au}
}

func (ar *AuthRouter) ConfigureAuthRouter(mux *mux.Router) {
	mux.HandleFunc("/signup", ar.SignupHandler).Methods("POST", "OPTIONS")
	mux.HandleFunc("/login", ar.LoginHandler).Methods("POST", "OPTIONS")
	mux.HandleFunc("/logout", ar.LogoutHandler).Methods("DELETE", "OPTIONS")
}
