package logger

import (
	"log"
	"os"
)

type Interface interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Logger struct {
	info  *log.Logger
	error *log.Logger
}

func New() *Logger {
	return &Logger{
		info:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		error: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.info.Printf(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.error.Printf(msg, args...)
}
