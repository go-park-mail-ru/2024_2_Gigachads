package user

import (
	"mail/api-service/pkg/utils"
	"net/http"
	"strings"
	"bytes"
	"time"
)

func (uh *UserRouter) GetAvatarHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)

	data, name, err := uh.UserUseCase.GetAvatar(email)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if strings.HasSuffix(name, ".png") {
		w.Header().Set("Content-Type", "image/png")
	} else {
		w.Header().Set("Content-Type", "image/jpeg")
	}

	
	w.WriteHeader(http.StatusOK)
	http.ServeContent(w, r, name, time.Now(), bytes.NewReader(data))
}
