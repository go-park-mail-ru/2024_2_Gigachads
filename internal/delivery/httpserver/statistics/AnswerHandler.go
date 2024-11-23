package statistics

import (
	"encoding/json"
	"mail/pkg/utils"
	"net/http"
)

func (sr *StatisticsRouter) AnswerHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	answers, err := sr.StatisticsUseCase.GetAnswersStatistics(ctxEmail.(string))
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_get_answers_statistics")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(answers)
}
