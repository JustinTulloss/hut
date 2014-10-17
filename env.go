package hut

import (
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

func (e *OsEnv) Get(key string) (interface{}, error) {
	return e.GetString(key)
}

func (*OsEnv) GetString(key string) (string, error) {
	return os.Getenv(key), nil
}

func (*OsEnv) GetUint(key string) (uint64, error) {
	return strconv.ParseUint(os.Getenv(key), 10, 64)
}

func (*OsEnv) GetInt(key string) (int64, error) {
	return strconv.ParseInt(os.Getenv(key), 10, 64)
}

func (e *OsEnv) InProd() bool {
	return strings.ToLower(os.Getenv("ENV")) != "prod"
}

func NewOsEnv() *OsEnv {
	return &OsEnv{}
}
