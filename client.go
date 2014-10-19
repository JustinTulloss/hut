package hut

import (
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func (s *Service) NewDbClient(driver string) (*sqlx.DB, error) {
	if db != nil {
		return db, nil
	}
	dbUrl, err := s.Env.GetString("database_url")
	if err != nil {
		return nil, err
	}
	conn, err := sqlx.Connect(driver, dbUrl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
