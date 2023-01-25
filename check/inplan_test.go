package check

import (
	"testing"

	"github.com/matt-FFFFFF/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/require"
)

func TestNumberOfResourcesInPlan(t *testing.T) {
	t.Parallel()

	tftest := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, tftest.Err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).NumberOfResourcesEquals(2).IfNotFail(t)
}

func TestNumberOfResourcesInPlanWithError(t *testing.T) {
	t.Parallel()

	tftest := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, tftest.Err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).NumberOfResourcesEquals(1).ErrorContains(t, "expected 1 resources, got 2")
}
