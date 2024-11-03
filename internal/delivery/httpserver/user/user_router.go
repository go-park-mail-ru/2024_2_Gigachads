package user

import (
	"github.com/gorilla/mux"
	"mail/internal/models"
)

type UserRouter struct {
	UserUseCase models.UserUseCase
}

func NewUserRouter(uu models.UserUseCase) *UserRouter {
	return &UserRouter{UserUseCase: uu}
}

func (ur *UserRouter) ConfigureUserRouter(mux *mux.Router) {
	mux.HandleFunc("/settings/avatar", ar.AvatarHandler).Methods("PUT", "OPTIONS")
	mux.HandleFunc("/settings/password", ar.PasswordHandler).Methods("PUT", "OPTIONS")
	mux.HandleFunc("/settings/name", ar.NameHandler).Methods("PUT", "OPTIONS")
}
