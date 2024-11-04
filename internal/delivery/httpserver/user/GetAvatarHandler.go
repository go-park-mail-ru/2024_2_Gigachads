package user

import (
	"mail/pkg/utils"
	"net/http"
)

func (uh *UserRouter) GetAvatarHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)

	buf, contType, err := uh.UserUseCase.GetAvatar(email)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", contType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+"/avatars/1.jpg"+"\"")

	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
