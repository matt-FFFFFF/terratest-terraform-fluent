package check

import (
	"fmt"
	"strings"
	"testing"
)

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
