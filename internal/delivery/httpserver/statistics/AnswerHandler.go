package statistics

import (
	"encoding/json"
	"mail/internal/models"
	"mail/pkg/utils"
	"net/http"
)

func (sr *StatisticsRouter) AnswerHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	var answers models.Answer
	err := json.NewDecoder(r.Body).Decode(&answers)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_request_body")
		return
	}

	err = sr.StatisticsUseCase.AnswersStatistics(answers.Action, answers.Value, ctxEmail.(string))
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_get_answers_statistics")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
