package email

import (
	"encoding/json"
	"mail/api-service/internal/models"
	"mail/api-service/pkg/utils"
	"net/http"
	"errors"
	"io"
)

func (er *EmailRouter) UploadAttachHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	// var reqfile models.File
	// err := json.NewDecoder(r.Body).Decode(&reqfile)
	// if err != nil {
	// 	utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
	// 	return
	// }

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

	file, _, err := r.FormFile("file")
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "error_with_file")
		return
	}
	defer file.Close()

	limitedReader := http.MaxBytesReader(w, file, 10*1024*1024)
	defer r.Body.Close()

	fileContent, err := io.ReadAll(limitedReader)
	if err != nil && !errors.Is(err, io.EOF) {
		if errors.As(err, new(*http.MaxBytesError)) {
			utils.ErrorResponse(w, r, http.StatusRequestEntityTooLarge, "too_big_body")
			return
		}
	}

	name := r.FormValue("name")

	path, err := er.EmailUseCase.UploadAttach(fileContent, name)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "error_with_upload_attach")
		return
	}
	jsonpath := models.FilePath{Path: path}

	result, err := json.Marshal(jsonpath)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "json_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
