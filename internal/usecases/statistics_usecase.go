package usecase

import (
	models "mail/internal/models"
)

type StatisticsService struct {
	StatisticsRepo models.StatisticsRepository
}

func NewStatisticsService(srepo models.StatisticsRepository) models.StatisticsUseCase {
	return &UserService{StatisticsRepo: srepo}
}

func (ss *StatisticsService) GetQuestionsStatistics() {

}
