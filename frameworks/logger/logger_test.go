package logger

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test-log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	logger, err := NewLogger("info", tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	testCases := []struct {
		level   string
		message string
	}{
		{"info", "Test info message"},
		{"warn", "Test warn message"},
		{"error", "Test error message"},
	}

	for _, tc := range testCases {
		switch tc.level {
		case "info":
			logger.Info(tc.message)
		case "warn":
			logger.Warn(tc.message)
		case "error":
			logger.Error(tc.message)
		}
	}

	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != len(testCases) {
		t.Fatalf("Expected %d log lines, got %d", len(testCases), len(lines))
	}

	for i, line := range lines {
		var logEntry map[string]interface{}
		err = json.Unmarshal([]byte(line), &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse log line: %v", err)
		}

		if logEntry["level"] != testCases[i].level {
			t.Errorf("Expected log level %s, got %s", testCases[i].level, logEntry["level"])
		}

		if logEntry["message"] != testCases[i].message {
			t.Errorf("Expected message '%s', got '%s'", testCases[i].message, logEntry["message"])
		}
	}
}
