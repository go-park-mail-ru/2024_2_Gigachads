package user

import (
	"net/http"
	"mail/pkg/utils"
)

func (uh *UserRouter) GetAvatarHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)
	
	buf, err := uh.UserUseCase.GetAvatar(email)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "multipart/form-data")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
