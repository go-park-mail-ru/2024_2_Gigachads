package middleware

import (
	"net/http"
	"mail/pkg/utils"
)

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				utils.ErrorResponse(w, r, http.StatusInternalServerError, "StatusServerError")
			}
		}()
		next.ServeHTTP(w, r)
	})
}