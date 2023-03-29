package setuptest

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"gopkg.in/matryer/try.v1"
)

type TestError struct {
	msg string
	//T   *testing.T
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
	if e == nil {
		return nil
	}
	return errors.New(e.msg)
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

// Retry is a configuration for retrying a terraform command.
type Retry struct {
	Max  int
	Wait time.Duration
}

// DefaultRetry is the default retry configuration.
// It will retry up to 5 times with a 1 minute wait between each attempt.
var DefaultRetry = Retry{
	Max:  5,
	Wait: time.Minute,
}

// DefaultRetry is the faster retry configuration.
// It will retry up to 6 times with a 20 second wait between each attempt.
var FastRetry = Retry{
	Max:  6,
	Wait: 20 * time.Second,
}

// DefaultRetry is the slower retry configuration.
// It will retry up to 20 times with a 1 minute wait between each attempt.
var SlowRetry = Retry{
	Max:  20,
	Wait: 1 * time.Minute,
}

// Apply runs terraform apply for the given Response and returns the error.
func (resp Response) Apply(t *testing.T) *TestError {

	_, err := terraform.ApplyE(t, resp.Options)
	return newTestError(err.Error())
}

// Apply runs terraform apply, then plan for the given Response and checks for any changes,
// it then returns the error.
func (resp Response) ApplyIdempotent(t *testing.T) *TestError {
	_, err := terraform.ApplyAndIdempotentE(t, resp.Options)
	return newTestError(err.Error())
}

// Apply runs terraform apply, then performs a retry loop with a plan.
// If the configuration is not idempotent, it will retry up to the specified number of times.
// It then returns the error.
func (resp Response) ApplyIdempotentRetry(t *testing.T, r Retry) *TestError {
	_, err := terraform.ApplyE(t, resp.Options)

	if err != nil {
		return newTestError(err.Error())
	}

	if try.MaxRetries < r.Max {
		try.MaxRetries = r.Max
	}

	err = try.Do(func(attempt int) (bool, error) {
		exitCode, err := terraform.PlanExitCodeE(t, resp.Options)
		if err != nil {
			t.Logf("terraform plan failed attempt %d/%d: waiting %s", attempt, r.Max, r.Wait)
			time.Sleep(r.Wait)
		}
		if exitCode != 0 {
			t.Logf("terraform not idempotent attempt %d/%d: waiting %s", attempt, r.Max, r.Wait)
			err = errors.New("terraform configuration not idempotent")
			time.Sleep(r.Wait)
		}
		return attempt < r.Max, err
	})

	if err != nil {
		return newTestError(err.Error())
	}

	return nil
}

// Destroy runs terraform destroy for the given Response and returns the error.
func (resp Response) Destroy(t *testing.T) *TestError {
	_, err := terraform.DestroyE(t, resp.Options)
	return newTestError(err.Error())
}

// DestroyWithRetry will retry the terraform destroy command up to the specified number of times.
func (resp Response) DestroyRetry(t *testing.T, r Retry) *TestError {
	if try.MaxRetries < r.Max {
		try.MaxRetries = r.Max
	}
	err := try.Do(func(attempt int) (bool, error) {
		_, err := terraform.DestroyE(t, resp.Options)
		if err != nil {
			t.Logf("terraform destroy failed attempt %d/%d: waiting %s", attempt, r.Max, r.Wait)
			time.Sleep(r.Wait)
		}
		return attempt < r.Max, err
	})
	if err != nil {
		return newTestErrorf("terraform destroy failed after %d attempts: %v", r.Max, err)
	}
	return nil
}
