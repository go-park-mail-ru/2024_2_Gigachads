package models

import (
	"time"
)

type Session struct {
	Name      string
	ID        string
	Time      time.Time
	UserLogin string
}
