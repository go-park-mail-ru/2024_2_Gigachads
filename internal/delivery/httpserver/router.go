package httpserver

import (
	"github.com/gorilla/mux"
	mw "mail/internal/delivery/middleware"
	usecase "mail/internal/usecases"
)

func ConfigureEmailRouter(privateMux *mux.Router, eu usecase.EmailUseCase) {
	emailRouter := NewEmailRouter(eu)
	privateMux.HandleFunc("/mail/inbox", emailRouter.InboxHandler).Methods("GET", "OPTIONS")
}

func ConfigureAuthRouter(publicMux *mux.Router, privateMux *mux.Router, uu usecase.UserUseCase) {
	authRouter := NewAuthRouter(uu)
	publicMux.HandleFunc("/signup", authRouter.SignupHandler).Methods("POST", "OPTIONS")
	publicMux.HandleFunc("/login", authRouter.LoginHandler).Methods("POST", "OPTIONS")
	privateMux.HandleFunc("/logout", authRouter.LogoutHandler).Methods("GET", "OPTIONS")
}

func ConfigureAuthMiddleware(privateMux *mux.Router, uu usecase.UserUseCase) {
	authMW := mw.NewAuthMW(uu)
	privateMux.Use(authMW.Handler)
}
