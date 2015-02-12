package hut

import (
	"os"
	"strconv"
	"strings"
)

type Env interface {
	Get(string) interface{}
	GetString(string) string
	GetUint(string) uint64
	GetInt(string) int64
	InProd() bool
}

// Reads the environment from the OS environment
type OsEnv struct{}

func getenv(key string) string {
	return os.Getenv(strings.ToUpper(key))
}

func (e *OsEnv) Get(key string) interface{} {
	return e.GetString(key)
}

func (*OsEnv) GetString(key string) string {
	return getenv(key)
}

func (*OsEnv) GetUint(key string) uint64 {
	val, err := strconv.ParseUint(getenv(key), 10, 64)
	if err != nil {
		panic("Could not parse " + key)
	}
	return val
}

func (*OsEnv) GetInt(key string) int64 {
	val, err := strconv.ParseInt(getenv(key), 10, 64)
	if err != nil {
		panic("Could not parse " + key)
	}
	return val
}

func (e *OsEnv) InProd() bool {
	return strings.ToLower(os.Getenv("ENV")) == "prod"
}

func NewOsEnv() *OsEnv {
	return &OsEnv{}
}

type MapEnv map[string]interface{}

func (e MapEnv) Get(key string) interface{} {
	return e[key]
}

func (e MapEnv) GetString(key string) string {
	s, ok := e[key].(string)
	if !ok {
		panic("Could not find " + key + " in env")
	}
	return s
}

func (e MapEnv) GetUint(key string) uint64 {
	n, ok := e[key].(uint64)
	if !ok {
		panic("Could not get " + key + " as a uint")
	}
	return n
}

func (e MapEnv) GetInt(key string) int64 {
	n, ok := e[key].(int64)
	if !ok {
		panic("Could not get " + key + " as an int")
	}
	return n
}

func (e MapEnv) InProd() bool {
	return e["ENV"].(string) == "prod"
}

func NewMapEnv() MapEnv {
	return MapEnv{}
}
