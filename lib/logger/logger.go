package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	FormatHTTPRequest(*http.Request) string
}

type MyLogger struct {
	logger *log.Logger
}

func NewLogger() *MyLogger {
	return &MyLogger{
		logger: log.New(os.Stdout, "", log.LUTC|log.Ldate|log.Lmicroseconds),
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

func (l *MyLogger) FormatHTTPRequest(r *http.Request) string {
	// TODO: include body
	headers := ""
	for key, values := range r.Header {
		for _, value := range values {
			headers = fmt.Sprintf("%s%s: %s\n", headers, key, value)
		}
	}

	return fmt.Sprintf(`HTTP Request
%s %s %s
Host: %s
%s`, r.Method, r.URL.Path, r.Proto, r.Host, headers)
}
