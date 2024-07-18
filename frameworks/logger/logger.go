package logger

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"
)

type Logger interface {
	Debug(msg string, fields ...LogField)
	Info(msg string, fields ...LogField)
	Warn(msg string, fields ...LogField)
	Error(msg string, fields ...LogField)
}

type LogField struct {
	Key   string
	Value interface{}
}

type simpleLogger struct {
	logger *log.Logger
}

func NewLogger(level string, filepath string) (Logger, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &simpleLogger{
		logger: log.New(file, "", 0), // Remove default prefixes
	}, nil
}

func (l *simpleLogger) log(level, msg string, fields ...LogField) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     strings.ToLower(level),
		"message":   msg,
	}

	for _, field := range fields {
		logEntry[field.Key] = field.Value
	}

	jsonLog, _ := json.Marshal(logEntry)
	l.logger.Println(string(jsonLog))
}

func (l *simpleLogger) Debug(msg string, fields ...LogField) {
	l.log("debug", msg, fields...)
}

func (l *simpleLogger) Info(msg string, fields ...LogField) {
	l.log("info", msg, fields...)
}

func (l *simpleLogger) Warn(msg string, fields ...LogField) {
	l.log("warn", msg, fields...)
}

func (l *simpleLogger) Error(msg string, fields ...LogField) {
	l.log("error", msg, fields...)
}

func Field(key string, value interface{}) LogField {
	return LogField{Key: key, Value: value}
}
