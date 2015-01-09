package hut

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Env interface {
	Get(string) (interface{}, error)
	GetString(string) (string, error)
	GetUint(string) (uint64, error)
	GetInt(string) (int64, error)
	InProd() bool
}

// Reads the environment from the OS environment
type OsEnv struct{}

func getenv(key string) string {
	return os.Getenv(strings.ToUpper(key))
}

func (e *OsEnv) Get(key string) (interface{}, error) {
	return e.GetString(key)
}

func (*OsEnv) GetString(key string) (string, error) {
	return getenv(key), nil
}

func (*OsEnv) GetUint(key string) (uint64, error) {
	return strconv.ParseUint(getenv(key), 10, 64)
}

func (*OsEnv) GetInt(key string) (int64, error) {
	return strconv.ParseInt(getenv(key), 10, 64)
}

func (e *OsEnv) InProd() bool {
	return strings.ToLower(os.Getenv("ENV")) == "prod"
}

func NewOsEnv() *OsEnv {
	return &OsEnv{}
}

type MapEnv map[string]interface{}

func (e MapEnv) Get(key string) (interface{}, error) {
	return e[key], nil
}

func (e MapEnv) GetString(key string) (string, error) {
	s, ok := e[key].(string)
	if !ok {
		return "", errors.New(key + " is not a string")
	}
	return s, nil
}

func (e MapEnv) GetUint(key string) (uint64, error) {
	n, ok := e[key].(uint64)
	if !ok {
		return 0, errors.New(key + " is not a uint64")
	}
	return n, nil
}

func (e MapEnv) GetInt(key string) (int64, error) {
	n, ok := e[key].(int64)
	if !ok {
		return 0, errors.New(key + " is not a int64")
	}
	return n, nil
}

func (e MapEnv) InProd() bool {
	return e["ENV"].(string) == "prod"
}

func NewMapEnv() MapEnv {
	return MapEnv{}
}
