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
		warn:  log.New(os.Stdout, "WARNING ", log.LstdFlags|log.Lshortfile),
		err:   log.New(os.Stdout, "ERROR ", log.LstdFlags|log.Lshortfile),
	}
}

// Shortcuts to make logging from a service a little less verbose
func (s *Service) Error(format string, v ...interface{}) {
	s.Log.Error().Printf(format, v...)
}

func (s *Service) Warning(format string, v ...interface{}) {
	s.Log.Warning().Printf(format, v...)
}

func (s *Service) Info(format string, v ...interface{}) {
	s.Log.Info().Printf(format, v...)
}

func (s *Service) Debug(format string, v ...interface{}) {
	s.Log.Debug().Printf(format, v...)
}
