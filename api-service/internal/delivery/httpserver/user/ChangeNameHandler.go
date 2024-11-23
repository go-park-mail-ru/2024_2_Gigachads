package user
import (
	"net/http"
	"encoding/json"
	"mail/api-service/pkg/utils"
	"mail/api-service/internal/models"
)

func (ar *UserRouter) ChangeNameHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)
	var nameModel models.ChangeName
	err := json.NewDecoder(r.Body).Decode(&nameModel)
	
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	nameModel.Name = utils.Sanitize(nameModel.Name)
		
	if !models.InputIsValid(nameModel.Name) {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_name")
		return
	}

	err = ar.UserUseCase.ChangeName(email, nameModel.Name)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
