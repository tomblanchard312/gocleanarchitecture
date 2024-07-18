package errors

import "fmt"

type AppError struct {
	Err     error
	Message string
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func New(message string) error {
	return &AppError{
		Message: message,
	}
}

func Wrap(err error, message string) error {
	return &AppError{
		Err:     err,
		Message: message,
	}
}
