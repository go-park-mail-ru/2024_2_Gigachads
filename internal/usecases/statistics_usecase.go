package usecase

import (

	// "context"
	// "fmt"
	// "net/http"
	models "mail/internal/models"
	// "mail/pkg/utils"
	// "os"	
)

type StatisticsService struct {
	StatisticsRepo    models.StatisticsRepository
}

func NewStatisticsService(srepo models.StatisticsRepository) models.StatisticsUseCase {
	return &StatisticsService{StatisticsRepo: srepo}
}

func (ss *StatisticsService) GetQuestionsStatistics(action string, email string) (models.Question, error) {
	return
}

func (ss *StatisticsService) AnswersStatistics(action string, value int, email string) error {
	return
}

func (ss *StatisticsService) GetStatistics() ([]models.Statistics, error) {
	return
}