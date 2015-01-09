package hut

import (
	"errors"
	"regexp"

	"github.com/JustinTulloss/firebase"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var firebaseishRe = regexp.MustCompile("^https://firebase")

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

func (s *Service) NewFirebaseClient() (*firebase.Client, error) {
	baseUrl, err := s.Env.GetString("FIREBASE_URL")
	if err != nil {
		return nil, err
	}
	if !firebaseishRe.MatchString(baseUrl) {
		return nil, errors.New(baseUrl + " does not appear to be a firebase url")
	}
	authToken, err := s.Env.GetString("FIREBASE_AUTH")
	if err != nil {
		return nil, err
	}
	return firebase.NewClient(baseUrl, authToken, nil), nil
}
