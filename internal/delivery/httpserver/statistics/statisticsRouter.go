package user

import (
	"mail/internal/models"

	"github.com/gorilla/mux"
)

type StatisticsRouter struct {
	StatisticsUseCase models.StatisticsUseCase
}

func NewStatisticsRouter(su models.StatisticsUseCase) *StatisticsRouter {
	return &StatisticsRouter{StatisticsUseCase: su}
}

func (sr *StatisticsRouter) ConfigureStatisticsRouter(mux *mux.Router) {
	mux.HandleFunc("/question", sr.QuestionHandler).Methods("GET", "OPTIONS")
	mux.HandleFunc("/answer", sr.AnswerHandler).Methods("POST", "OPTIONS")
	mux.HandleFunc("/statistics", sr.StatisticsHandler).Methods("GET", "OPTIONS")
}