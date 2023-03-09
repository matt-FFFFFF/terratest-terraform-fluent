package setuptest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/depth1", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	require.NoError(t, test.Apply(t))
	require.NoError(t, test.Destroy(t))
}
