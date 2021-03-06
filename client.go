package hut

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/JustinTulloss/firebase"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var firebaseishRe = regexp.MustCompile(`^https://.*\.firebaseio`)

func (s *Service) NewDbClient(driver string) (*sqlx.DB, error) {
	if db != nil {
		return db, nil
	}
	dbUrl := s.Env.GetString("database_url")
	conn, err := sqlx.Connect(driver, dbUrl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (s *Service) NewFirebaseClient() (firebase.Client, error) {
	baseUrl := s.Env.GetString("FIREBASE_URL")
	if !firebaseishRe.MatchString(baseUrl) {
		return nil, errors.New(baseUrl + " does not appear to be a firebase url")
	}
	authToken := s.Env.GetString("FIREBASE_AUTH")
	return firebase.NewClient(baseUrl, authToken, nil), nil
}

func (s *Service) NewRedisClient() (redis.Conn, error) {
	redisTcpAddr := s.Env.GetString("REDIS_PORT_6379_TCP_ADDR")
	redisPort := s.Env.GetString("REDIS_PORT_6379_TCP_PORT")
	redisAddress := fmt.Sprintf(
		"%s:%s",
		redisTcpAddr,
		redisPort,
	)
	c, err := redis.Dial("tcp", redisAddress)
	if err != nil {
		return nil, err
	}
	return c, nil
}
