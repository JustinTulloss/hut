package hut

import (
	"fmt"
	"errors"
	"regexp"

	"github.com/JustinTulloss/firebase"
	"github.com/jmoiron/sqlx"
	"github.com/garyburd/redigo/redis"
)

var db *sqlx.DB
var firebaseishRe = regexp.MustCompile(`^https://.*\.firebaseio`)

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

func (s *Service) NewFirebaseClient() (firebase.Client, error) {
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


func (s *Service) NewRedisClient() (redis.Conn, *redis.PubSubConn, error) {
	redisTcpAddr, err := s.Env.GetString("REDIS_PORT_6379_TCP_ADDR")
	if err != nil {
		return nil, nil, err
	}
	redisPort, err := s.Env.GetString("REDIS_PORT_6379_TCP_PORT")
	if err != nil {
		return nil, nil, err
	}
	redisAddress := fmt.Sprintf(
		"%s:%s",
		redisTcpAddr,
		redisPort,
	)
	c, err := redis.Dial("tcp", redisAddress)
	if err != nil {
		return nil, nil, err
	}
	subscribeConn := &redis.PubSubConn{c}

	publishConn, err := redis.Dial("tcp", redisAddress)
	if err != nil {
		return nil, nil, err
	}
	return publishConn, subscribeConn, err
}

