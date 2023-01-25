package setuptest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirs(t *testing.T) {
	t.Parallel()

	test := Dirs("testdata/depth1", "").WithVars(map[string]interface{}{}).InitPlanShow(t)
	require.NoError(t, test.Err)
}
