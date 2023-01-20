package check

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// CheckError is a simple error type that allows us to chain checks using methods.
// However due to Go's way of handling interface types, when a nil *CheckError value is used as
// an error type (e.g. when passed into a func accepting an error) the underlying concrete
// type is *CheckError and it will not pass the usual error != nil check.
// Instead use the AsError() method to get a regular error type.
// Or use the reflect package `val := reflect.ValueOf(myCheckError); val.IsNil()`.
type CheckError struct {
	msg string
}

func newCheckError(msg string) *CheckError {
	return &CheckError{
		msg: msg,
	}
}

func newCheckErrorf(format string, args ...any) *CheckError {
	return &CheckError{
		msg: fmt.Sprintf(format, args...),
	}
}

// Implement Error interface
func (e *CheckError) Error() string {
	return e.msg
}

// AsError returns a regular error type that can be used in the usual way.
func (e *CheckError) AsError() error {
	if e == nil {
		return nil
	}
	return errors.New(e.msg)
}

func (e *CheckError) IfNotFail(t *testing.T) {
	if e != nil {
		t.Errorf(e.msg)
	}
}

func (e *CheckError) IfNotFailNow(t *testing.T) {
	if e != nil {
		t.Fatalf(e.msg)
	}
}

func (e *CheckError) ErrorContains(t *testing.T, substr string) {
	if !strings.Contains(e.msg, substr) {
		t.Errorf("error '%s' does not contain substring '%s'", e.msg, substr)
	}
}

func (e *CheckError) ErrorNotContains(t *testing.T, substr string) {
	if strings.Contains(e.msg, substr) {
		t.Errorf("error '%s' does contain substring '%s'", e.msg, substr)
	}
}
