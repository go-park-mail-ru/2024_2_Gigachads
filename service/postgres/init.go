package postgres

import (
	"database/sql"
	"fmt"
	"mail/config"

	_ "github.com/lib/pq"
)

func Init(cfg *config.Config) (*sql.DB, error) {
	dbConfig := cfg.Postgres
	dataConnection := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		dbConfig.IP, dbConfig.Port, dbConfig.DBname, dbConfig.User, dbConfig.Password)

	db, err := sql.Open("postgres", dataConnection)
	if err != nil {
		return nil, err
	}

	return db, nil
}
