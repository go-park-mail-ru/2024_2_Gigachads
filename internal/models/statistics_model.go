package models

type GetQuestion struct {
	Action string `json:"action"`
}

type Question struct {
	Description string `json:"description"`
	Type        string `json:"type"`
}

type Answer struct {
	Action string `json:"action"`
	Value  int    `json:"value"`
}

type Statistics struct {
	Action  string  `json:"action"`
	Amount  int     `json:"amount"`
	Average float32 `json:"average"`
}

type StatisticsUseCase interface {
	GetQuestionsStatistics(action string) (Question, error)
	AnswersStatistics(action string, value int, email string) error
	GetStatistics() ([]Statistics, error)
}

type StatisticsRepository interface {
	GetQuestionsStatistics(action string) (Question, error)
	AnswersStatistics(action string, value int, email string) error
	GetStatistics() ([]Statistics, error)
}
