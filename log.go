package hut

import (
	"log"
	"os"

	logxi "github.com/mgutz/logxi/v1"
)

type Logger interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{}) error
	Error(msg string, args ...interface{}) error
	Fatal(msg string, args ...interface{})
	Log(level int, msg string, args []interface{})

	SetLevel(int)
	IsTrace() bool
	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
	// Error, Fatal not needed, those SHOULD always be logged
}

type StdLog struct {
	trace *log.Logger
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
}

func (l *StdLog) Trace(msg string, args ...interface{}) {
	l.trace.Printf(msg, args...)
}
func (l *StdLog) Debug(msg string, args ...interface{}) {
	l.debug.Printf(msg, args...)
}
func (l *StdLog) Info(msg string, args ...interface{}) {
	l.info.Printf(msg, args...)
}
func (l *StdLog) Warn(msg string, args ...interface{}) error {
	l.warn.Printf(msg, args...)
	return nil
}
func (l *StdLog) Error(msg string, args ...interface{}) error {
	l.err.Printf(msg, args...)
	return nil
}
func (l *StdLog) Fatal(msg string, args ...interface{}) {
	l.err.Fatalf(msg, args...)
}

func NewStdLog() *StdLog {
	return &StdLog{
		trace: log.New(os.Stdout, "TRACE ", 0),
		debug: log.New(os.Stdout, "DEBUG ", 0),
		info:  log.New(os.Stdout, "INFO ", 0),
		warn:  log.New(os.Stdout, "WARNING ", log.Lshortfile),
		err:   log.New(os.Stdout, "ERROR ", log.Lshortfile),
	}
}

// Shortcuts to make logging from a service a little less verbose
func (s *Service) Trace(msg string, args ...interface{}) {
	s.Log.Trace(msg, args...)
}
func (s *Service) Debug(msg string, args ...interface{}) {
	s.Log.Debug(msg, args...)
}
func (s *Service) Info(msg string, args ...interface{}) {
	s.Log.Info(msg, args...)
}
func (s *Service) Warn(msg string, args ...interface{}) error {
	return s.Log.Warn(msg, args...)
}
func (s *Service) Error(msg string, args ...interface{}) error {
	return s.Log.Error(msg, args...)
}
func (s *Service) Fatal(msg string, args ...interface{}) {
	s.Log.Fatal(msg, args...)
}

func NewLogxiLog(name string) logxi.Logger {
	return logxi.New(name)
}
