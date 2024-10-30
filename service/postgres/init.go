package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"mail/config"
)

func Init(cfg *config.Config) (*sql.DB, error) {
	dbConfig := cfg.Postgres
	dataConnection := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
		dbConfig.IP, dbConfig.Port, dbConfig.DBname, dbConfig.User, dbConfig.Password)

	db, err := sql.Open("postgres", dataConnection)
	if err != nil {
		return nil, err
	}

	return db, nil
}
