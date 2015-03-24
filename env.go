package hut

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Env interface {
	Get(string) interface{}
	GetString(string) string
	GetUint(string) uint64
	GetInt(string) int64
	GetBool(string) bool
	GetTCPServiceAddress(string, int) string
	GetUDPServiceAddress(string, int) string
	MustGetTCPServiceAddress(string, int) string
	MustGetUDPServiceAddress(string, int) string
	InProd() bool
}

// Reads the environment from the OS environment
type OsEnv struct{}

// Given a docker link like name, pull the connection details out of the
// environment. https://docs.docker.com/userguide/dockerlinks/#environment-variables
func makeDockerLinkAddress(env Env, serviceName string, port int, protocol string) string {
	prefix := fmt.Sprintf("%s_PORT_%d_%s",
		strings.ToUpper(serviceName), port, protocol)

	envAddr := env.GetString(fmt.Sprintf("%s_ADDR", prefix))
	envPort := env.GetString(fmt.Sprintf("%s_PORT", prefix))
	// The min information we need is a port
	if envPort == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s", envAddr, envPort)
}
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

func (*OsEnv) GetBool(key string) bool {
	if getenv(key) == "" {
		return false
	}
	val, err := strconv.ParseBool(getenv(key))
	if err != nil {
		panic("Could not parse " + key)
	}
	return val
}

func (e *OsEnv) GetTCPServiceAddress(name string, port int) string {
	return makeDockerLinkAddress(e, name, port, "TCP")
}

func (e *OsEnv) GetUDPServiceAddress(name string, port int) string {
	return makeDockerLinkAddress(e, name, port, "UDP")
}

func (e *OsEnv) MustGetTCPServiceAddress(name string, port int) string {
	addr := e.GetTCPServiceAddress(name, port)
	if addr == "" {
		panic(fmt.Sprintf("Environment for %s was not configured", name))
	}
	return addr
}

func (e *OsEnv) MustGetUDPServiceAddress(name string, port int) string {
	addr := e.GetUDPServiceAddress(name, port)
	if addr == "" {
		panic(fmt.Sprintf("Environment for %s was not configured", name))
	}
	return addr
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

func (e MapEnv) GetBool(key string) bool {
	b, ok := e[key].(bool)
	if !ok {
		return false
	}
	return b
}

func (e MapEnv) GetTCPServiceAddress(name string, port int) string {
	return makeDockerLinkAddress(e, name, port, "TCP")
}

func (e MapEnv) GetUDPServiceAddress(name string, port int) string {
	return makeDockerLinkAddress(e, name, port, "UDP")
}

func (e MapEnv) MustGetTCPServiceAddress(name string, port int) string {
	addr := e.GetTCPServiceAddress(name, port)
	if addr == "" {
		panic(fmt.Sprintf("Environment for %s was not configured", name))
	}
	return addr
}

func (e MapEnv) MustGetUDPServiceAddress(name string, port int) string {
	addr := e.GetUDPServiceAddress(name, port)
	if addr == "" {
		panic(fmt.Sprintf("Environment for %s was not configured", name))
	}
	return addr
}

func (e MapEnv) InProd() bool {
	return e["ENV"].(string) == "prod"
}

func NewMapEnv() MapEnv {
	return MapEnv{}
}
