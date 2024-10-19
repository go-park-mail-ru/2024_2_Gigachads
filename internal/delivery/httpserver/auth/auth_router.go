package auth

import (
	"github.com/gorilla/mux"
	"mail/internal/models"
)

type AuthRouter struct {
	UserUseCase models.UserUseCase
}

func NewAuthRouter(uu models.UserUseCase) *AuthRouter {
	return &AuthRouter{UserUseCase: uu}
}

func (ar *AuthRouter) ConfigureAuthRouter(publicMux *mux.Router, privateMux *mux.Router) {
	publicMux.HandleFunc("/signup", ar.SignupHandler).Methods("POST", "OPTIONS")
	publicMux.HandleFunc("/login", ar.LoginHandler).Methods("POST", "OPTIONS")
	privateMux.HandleFunc("/logout", ar.LogoutHandler).Methods("GET", "OPTIONS")
}
