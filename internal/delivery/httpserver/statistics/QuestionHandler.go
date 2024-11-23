package statistics

import (
	"encoding/json"
	"mail/internal/models"
	"mail/pkg/utils"
	"net/http"
)

func (sr *StatisticsRouter) QuestionHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	var questions models.GetQuestion
	err := json.NewDecoder(r.Body).Decode(&questions)

	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_request_body")
		return
	}

	question, err := sr.StatisticsUseCase.GetQuestionsStatistics(questions.Action, ctxEmail.(string))
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_get_questions_statistics")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(question)
}
