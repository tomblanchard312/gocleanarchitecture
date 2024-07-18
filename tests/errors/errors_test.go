package errors

import (
	"gocleanarchitecture/errors"
	"testing"
)

func TestNew(t *testing.T) {
	msg := "test error"
	err := errors.New(msg)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != msg {
		t.Errorf("Expected error message '%s', got '%s'", msg, err.Error())
	}
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	msg := "wrapped error"
	wrappedErr := errors.Wrap(originalErr, msg)

	if wrappedErr == nil {
		t.Fatal("Expected wrapped error, got nil")
	}

	expectedMsg := msg + ": " + originalErr.Error()
	if wrappedErr.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, wrappedErr.Error())
	}

	appErr, ok := wrappedErr.(*errors.AppError)
	if !ok {
		t.Fatal("Expected *AppError type")
	}

	if appErr.Err != originalErr {
		t.Error("Original error not preserved in wrapped error")
	}
}
