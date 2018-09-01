package logger

import (
	"log"
	"os"
)

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}

type MyLogger struct {
	logger *log.Logger
}

func NewLogger() *MyLogger {
	return &MyLogger{
		logger: log.New(os.Stdout, "", log.Lmicroseconds),
	}
}

func (l *MyLogger) Print(v ...interface{}) {
	l.logger.SetOutput(os.Stdout)
	l.logger.Print(v...)
}

func (l *MyLogger) Printf(format string, v ...interface{}) {
	l.logger.SetOutput(os.Stdout)
	l.logger.Printf(format, v...)
}

func (l *MyLogger) Error(v ...interface{}) {
	l.logger.SetOutput(os.Stderr)
	l.logger.Print(v...)
}

func (l *MyLogger) Errorf(format string, v ...interface{}) {
	l.logger.SetOutput(os.Stderr)
	l.logger.Printf(format, v...)
}
