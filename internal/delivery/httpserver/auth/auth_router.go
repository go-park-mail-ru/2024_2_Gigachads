package auth

import (
	"mail/internal/models"

	"github.com/gorilla/mux"
)

type AuthRouter struct {
	UserUseCase models.UserUseCase
}

func NewAuthRouter(uu models.UserUseCase) *AuthRouter {
	return &AuthRouter{UserUseCase: uu}
}

func (ar *AuthRouter) ConfigureAuthRouter(mux *mux.Router) {
	mux.HandleFunc("/signup", ar.SignupHandler).Methods("POST", "OPTIONS")
	mux.HandleFunc("/login", ar.LoginHandler).Methods("POST", "OPTIONS")
	mux.HandleFunc("/logout", ar.LogoutHandler).Methods("DELETE", "OPTIONS")
}
