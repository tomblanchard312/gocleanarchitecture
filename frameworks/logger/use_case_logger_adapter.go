package logger

import "gocleanarchitecture/usecases"

// Adapter to make framework logger compatible with use case logger interface
type UseCaseLoggerAdapter struct {
	logger Logger
}

func NewUseCaseLoggerAdapter(logger Logger) usecases.Logger {
	return &UseCaseLoggerAdapter{logger: logger}
}

func (a *UseCaseLoggerAdapter) Error(msg string, fields ...interface{}) {
	// Convert variadic interface{} to LogField format
	var logFields []LogField
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				logFields = append(logFields, LogField{Key: key, Value: fields[i+1]})
			}
		}
	}
	a.logger.Error(msg, logFields...)
}
