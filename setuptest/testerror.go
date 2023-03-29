package setuptest

import (
	"errors"
	"fmt"
	"testing"
)

type TestError struct {
	msg string
}

func newTestError(msg string) *TestError {
	return &TestError{
		msg: msg,
	}
}

func newTestErrorf(format string, args ...any) *TestError {
	return &TestError{
		msg: fmt.Sprintf(format, args...),
	}
}

// Implement Error interface
func (e *TestError) Error() string {
	return e.msg
}

// AsError returns a regular error type that can be used in the usual way.
func (e *TestError) AsError() error {
	if e != nil || e.msg != "" {
		return errors.New(e.msg)
	}
	return nil
}

func (e *TestError) IfNotFail(t *testing.T) {
	if e != nil {
		t.Errorf(e.msg)
	}
}

func (e *TestError) IfNotFailNow(t *testing.T) {
	if e != nil {
		t.Fatalf(e.msg)
	}
}
