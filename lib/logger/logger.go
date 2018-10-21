package logger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
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
	mu     *sync.Mutex
	logger *log.Logger
}

func NewLogger() *MyLogger {
	return &MyLogger{
		mu:     new(sync.Mutex),
		logger: log.New(os.Stdout, "", log.LUTC|log.Ldate|log.Lmicroseconds),
	}
}

func (l *MyLogger) Print(v ...interface{}) {
	l.mu.Lock()
	l.logger.SetOutput(os.Stdout)
	l.logger.Print(v...)
	l.mu.Unlock()
}

func (l *MyLogger) Printf(format string, v ...interface{}) {
	l.mu.Lock()
	l.logger.SetOutput(os.Stdout)
	l.logger.Printf(format, v...)
	l.mu.Unlock()
}

func (l *MyLogger) Error(v ...interface{}) {
	l.mu.Lock()
	l.logger.SetOutput(os.Stderr)
	l.logger.Print(v...)
	l.mu.Unlock()
}

func (l *MyLogger) Errorf(format string, v ...interface{}) {
	l.mu.Lock()
	l.logger.SetOutput(os.Stderr)
	l.logger.Printf(format, v...)
	l.mu.Unlock()
}

func (l *MyLogger) FormatHTTPRequest(r *http.Request) string {
	log := fmt.Sprintf(`HTTP Request
%s %s %s
Host: %s
%s`, r.Method, r.URL.Path, r.Proto, r.Host, formatHeaders(r.Header))

	return readAndAppendBody(l, r.Body, log)
}

func (l *MyLogger) FormatHTTPResponse(r *http.Response) string {
	log := fmt.Sprintf(`HTTP Response
%s
%s %s %s
Host: %s
%s`, r.Status, r.Request.Method, r.Request.URL.Path, r.Proto, r.Request.Host, formatHeaders(r.Header))

	return readAndAppendBody(l, r.Body, log)
}

func formatHeaders(header map[string][]string) string {
	headers := ""
	for key, values := range header {
		for _, value := range values {
			headers = fmt.Sprintf("%s%s: %s\n", headers, key, value)
		}
	}
	return headers
}

func readAndAppendBody(logger *MyLogger, body io.ReadCloser, log string) string {
	if body == nil {
		return log
	}

	data, err := ioutil.ReadAll(body)
	if err != nil {
		logger.Errorf("error reading request body: %s\nrequest: %s", err.Error(), log)
		return log
	}

	closeErr := body.Close()
	if closeErr != nil {
		logger.Errorf("error closing request body: %s\nrequest: %s", closeErr.Error(), log)
	}

	body = ioutil.NopCloser(bytes.NewBuffer(data))
	return fmt.Sprintf("%s%s", log, data)
}
