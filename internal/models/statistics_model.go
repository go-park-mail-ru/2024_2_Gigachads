package models

import (
)

type GetQuestion struct {
	Action     string `json:"action"`
}

type Question struct {
	Description     string `json:"description"`
	Type    		string `json:"type"`
}

type Statistics struct {
	Action     string `json:"action"`
	Amount     int 	  `json:"amount"`
	Average    float32  `json:"average"`
}

type StatisticsUseCase interface {
	GetQuestionsStatistics()
	GetAnswersStatistics()
	GetStatistics()
}

type StatisticsRepository interface {
	GetQuestionsStatistics(signup *User) (*User, error)
	GetAnswersStatistics()
	GetStatistics()
}