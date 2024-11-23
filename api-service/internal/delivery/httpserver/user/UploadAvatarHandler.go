package user

import (
	"mail/api-service/pkg/utils"
	"net/http"
	"io"
	"errors"
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
	defer func() {
		if err := r.MultipartForm.RemoveAll(); err != nil {
			utils.ErrorResponse(w, r, http.StatusInternalServerError, "error_removing_temp_files")
		}
	}()

	file, _, err := r.FormFile("avatar")
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "error_with_file")
		return
	}
	defer file.Close()

	limitedReader := http.MaxBytesReader(w, file, 10 * 1024 * 1024)
	defer r.Body.Close()

	fileContent, err := io.ReadAll(limitedReader)
	if err != nil && !errors.Is(err, io.EOF) {
		if errors.As(err, new(*http.MaxBytesError)) {
			utils.ErrorResponse(w, r, http.StatusRequestEntityTooLarge, "too_big_body")
			return
		}
	}

	err = uh.UserUseCase.ChangeAvatar(fileContent, email)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
