package hut

import (
	"log"
	"os"
)

type Log interface {
	Debug() *log.Logger
	Info() *log.Logger
	Warning() *log.Logger
	Error() *log.Logger
}

type StdLog struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
}

func (l *StdLog) Debug() *log.Logger {
	return l.debug
}

func (l *StdLog) Info() *log.Logger {
	return l.info
}

func (l *StdLog) Warning() *log.Logger {
	return l.warn
}

func (l *StdLog) Error() *log.Logger {
	return l.err
}

func NewStdLog() *StdLog {
	return &StdLog{
		debug: log.New(os.Stdout, "DEBUG ", log.LstdFlags),
		info:  log.New(os.Stdout, "INFO ", log.LstdFlags),
		warn:  log.New(os.Stdout, "WARNING ", log.LstdFlags),
		err:   log.New(os.Stdout, "ERROR ", log.LstdFlags),
	}
}
