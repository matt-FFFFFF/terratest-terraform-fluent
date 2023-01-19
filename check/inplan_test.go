package check

import (
	"testing"

	"github.com/matt-FFFFFF/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumberOfResourcesInPlan(t *testing.T) {
	t.Parallel()

	tftest := setuptest.Dirs(basicTestData, "").WithVars(nil).InitAndPlanAndShowWithStruct(t)
	require.NoError(t, tftest.Err)
	defer tftest.Cleanup()
	assert.NoError(
		t,
		InPlan(tftest.Plan).NumberOfResourcesEquals(2),
	)
}

func TestNumberOfResourcesInPlanWithError(t *testing.T) {
	t.Parallel()

	tftest := setuptest.Dirs(basicTestData, "").WithVars(nil).InitAndPlanAndShowWithStruct(t)
	require.NoError(t, tftest.Err)
	defer tftest.Cleanup()
	assert.ErrorContains(
		t,
		InPlan(tftest.Plan).NumberOfResourcesEquals(1),
		"expected 1 resources, got 2",
	)
}
