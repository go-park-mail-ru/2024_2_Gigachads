package user

import (
	"net/http"
	"mail/pkg/utils"
)


func (uh *UserRouter) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)
	
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "error_with_parsing_file")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "error_with_file")
		return
	}
	defer file.Close()
	
	err = uh.UserUseCase.ChangeAvatar(file, *header, email)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "error_with_saving_avatar")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
