package setuptest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/depth1", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	test.Apply(t).IfNotFail(t)
	test.Destroy(t).IfNotFail(t)
}

func TestApplyIdempotent(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/depth1", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	test.ApplyIdempotent(t).IfNotFail(t)
	test.Destroy(t).IfNotFail(t)
}
