package models

type Error struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}
