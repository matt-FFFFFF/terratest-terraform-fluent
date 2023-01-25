package setuptest

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"gopkg.in/matryer/try.v1"
)

type Retry struct {
	Max  int
	Wait time.Duration
}

var DefaultRetry = Retry{
	Max:  5,
	Wait: time.Minute,
}

var FastRetry = Retry{
	Max:  6,
	Wait: 20 * time.Second,
}

var LongRetry = Retry{
	Max:  20,
	Wait: 1 * time.Minute,
}

func (resp SetupTestResponse) Apply(t *testing.T) error {
	_, err := terraform.ApplyE(t, resp.Options)
	return err
}

func (resp SetupTestResponse) ApplyIdempotent(t *testing.T) error {
	_, err := terraform.ApplyAndIdempotentE(t, resp.Options)
	return err
}

func (resp SetupTestResponse) Destroy(t *testing.T) error {
	_, err := terraform.DestroyE(t, resp.Options)
	return err
}

func (resp SetupTestResponse) DestroyWithRetry(t *testing.T, r Retry) {
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
	}
}
