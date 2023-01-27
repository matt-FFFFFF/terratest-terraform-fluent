package setuptest

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"gopkg.in/matryer/try.v1"
)

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
func (resp Response) Apply(t *testing.T) error {
	_, err := terraform.ApplyE(t, resp.Options)
	return err
}

// Apply runs terraform apply, then plan for the given Response and checks for any changes,
// it then returns the error.
func (resp Response) ApplyIdempotent(t *testing.T) error {
	_, err := terraform.ApplyAndIdempotentE(t, resp.Options)
	return err
}

// Destroy runs terraform destroy for the given Response and returns the error.
func (resp Response) Destroy(t *testing.T) error {
	_, err := terraform.DestroyE(t, resp.Options)
	return err
}

// DestroyWithRetry will retry the terraform destroy command up to the specified number of times.
func (resp Response) DestroyWithRetry(t *testing.T, r Retry) error {
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
		t.Logf("terraform destroy failed after %d attempts: %v", r.Max, err)
		return err
	}
	return nil
}
