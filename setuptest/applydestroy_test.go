package setuptest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/depth1", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	test.Apply(t).ErrorIsNil(t)
	test.Destroy(t).ErrorIsNil(t)
}

func TestApplyIdempotent(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/depth1", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	test.ApplyIdempotent(t).ErrorIsNil(t)
	test.Destroy(t).ErrorIsNil(t)
}

// TestApplyRetryIdempotentFail is a test that will retry the apply if it fails.
// We simulate this with a resource that has a property using `timestamp()` which
// will always be different on every plan.
func TestApplyIdempotentRetryFail(t *testing.T) {
	t.Parallel()

	rty := Retry{
		Max:  3,
		Wait: time.Second * 20,
	}
	test, err := Dirs("testdata/applyidempotentretryfail", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer test.Destroy(t).ErrorIsNil(t)
	err = test.ApplyIdempotentRetry(t, rty).AsError()
	assert.ErrorContains(t, err, "not idempotent")
}
