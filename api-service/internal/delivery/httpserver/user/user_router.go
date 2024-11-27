package user

import (
	"github.com/gorilla/mux"
	"mail/api-service/internal/models"
)

type UserRouter struct {
	UserUseCase models.UserUseCase
}

func NewUserRouter(uu models.UserUseCase) *UserRouter {
	return &UserRouter{UserUseCase: uu}
}

func (ur *UserRouter) ConfigureUserRouter(mux *mux.Router) {
	mux.HandleFunc("/settings/avatar", ur.UploadAvatarHandler).Methods("PUT", "OPTIONS")
	mux.HandleFunc("/settings/password", ur.ChangePasswordHandler).Methods("PUT", "OPTIONS")
	mux.HandleFunc("/settings/name", ur.ChangeNameHandler).Methods("PUT", "OPTIONS")
	mux.HandleFunc("/settings/avatar", ur.GetAvatarHandler).Methods("GET", "OPTIONS")
}
