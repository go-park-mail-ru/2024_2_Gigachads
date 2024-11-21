package middleware

import (
	"github.com/gorilla/mux"
	"mail/config"
	"net/http"
)

func ConfigureMWs(cfg *config.Config, mux *mux.Router, authMW *AuthMiddleware) http.Handler {
	handler := authMW.Handler(mux)
	handler = CORS(handler, cfg)
	return handler
}
