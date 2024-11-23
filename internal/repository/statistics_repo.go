package repository

import (
	"database/sql"
	"mail/internal/models"
	"mail/pkg/logger"
)

type StatisticsRepositoryService struct {
	repo   *sql.DB
	logger logger.Logable
}

func NewStatisticsRepositoryService(db *sql.DB, l logger.Logable) *StatisticsRepositoryService {
	return &StatisticsRepositoryService{repo: db, logger: l}
}

func (sr *StatisticsRepositoryService) GetQuestionsStatistics(action string) (models.Question, error) {
	row := sr.repo.QueryRow(
		`SELECT description, type FROM question WHERE action = $1`, action)
	question := models.Question{}
	err := row.Scan(&question.Description, &question.Type)
	if err != nil {
		sr.logger.Error(err.Error())
		return models.Question{}, err
	}
	return question, nil
}

func (sr *StatisticsRepositoryService) AnswersStatistics(action string, value int, email string) error {
	_, err := sr.repo.Exec(
		`INSERT INTO "answer" (action, value, user_email) VALUES ($1, $2, $3)`, action, value, email)
	if err != nil {
		sr.logger.Error(err.Error())
		return err
	}
	return nil
}

func (sr *StatisticsRepositoryService) GetStatistics() ([]models.Statistics, error) {
	rows, err := sr.repo.Query(
		`SELECT AVERAGE(value), COUNT(value), action FROM "answer" GROUP BY action`)
	if err != nil {
		sr.logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()
	res := make([]models.Statistics, 0)
	for rows.Next() {
		stat := models.Statistics{}
		err := rows.Scan(
			&stat.Average,
			&stat.Amount,
			&stat.Action,
		)
		if err != nil {
			sr.logger.Error(err.Error())
			return nil, err
		}
		res = append(res, stat)
	}
	return res, nil
}