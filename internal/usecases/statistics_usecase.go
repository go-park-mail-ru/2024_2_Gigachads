package usecase

import (

	"context"
	"fmt"
	"net/http"
	models "mail/internal/models"
	"mail/pkg/utils"
	"os"	
)

type StatisticsService struct {
	StatisticsRepo    models.StatisticsRepository
}

func NewStatisticsService(srepo models.StatisticsRepository) models.StatisticsUseCase {
	return &UserService{StatisticsRepo: srepo}
}

func (ss *StatisticsService) GetQuestionsStatistics() {
	
}