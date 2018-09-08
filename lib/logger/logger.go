package logger

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	FormatHTTPResponse(*http.Response) string
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
	log := fmt.Sprintf(`HTTP Request
%s %s %s
Host: %s
%s`, r.Method, r.URL.Path, r.Proto, r.Host, formatHeader(r.Header))

	if r.Body != nil {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			l.Errorf(err.Error(), log)
		} else {
			r.Body.Close()
			r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
			log = fmt.Sprintf("%s%s", log, data)
		}
	}

	return log
}

func (l *MyLogger) FormatHTTPResponse(r *http.Response) string {
	log := fmt.Sprintf(`HTTP Response
%s
%s %s %s
Host: %s
%s`, r.Status, r.Request.Method, r.Request.URL.Path, r.Proto, r.Request.Host, formatHeader(r.Header))

	if r.Body != nil {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			l.Errorf(err.Error(), log)
		} else {
			r.Body.Close()
			r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
			log = fmt.Sprintf("%s%s", log, data)
		}
	}

	return log
}

func formatHeader(header map[string][]string) string {
	headers := ""
	for key, values := range header {
		for _, value := range values {
			headers = fmt.Sprintf("%s%s: %s\n", headers, key, value)
		}
	}
	return headers
}
