package auth

import (
	"mail/api-service/pkg/utils"
	"net/http"
)

func (ar *AuthRouter) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)

	err := ar.AuthUseCase.Logout(r.Context(), email)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "error_with_logout")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
