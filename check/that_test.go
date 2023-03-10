package check

import (
	"testing"

	"github.com/matt-FFFFFF/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	basicTestData = "testdata/basic"
)

func TestHasValueStrings(t *testing.T) {
	t.Parallel()

	tftest := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, tftest.Err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test").Key("content").HasValue("test").IfNotFail(t)
}

func TestHasValueStringsNotEqualError(t *testing.T) {
	t.Parallel()

	tftest := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, tftest.Err)
	defer tftest.Cleanup()
	assert.ErrorContains(
		t,
		InPlan(tftest.Plan).That("local_file.test").Key("content").HasValue("throwError"),
		"attribute content, planned value test not equal to assertion throwError",
	)
}

func TestHasValueStringsToInt(t *testing.T) {
	t.Parallel()

	tftest := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, tftest.Err)
	defer tftest.Cleanup()
	assert.Error(
		t,
		InPlan(tftest.Plan).That("local_file.test_int").Key("content").HasValue(123),
	)
}

func TestKeyNotExistsError(t *testing.T) {
	t.Parallel()

	tftest := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	defer tftest.Cleanup()
	assert.ErrorContains(
		t,
		InPlan(tftest.Plan).That("local_file.test").Key("not_exists").Exists(),
		"key not_exists not found in resource",
	)
}

func TestKeyNotExists(t *testing.T) {
	t.Parallel()

	tftest := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	defer tftest.Cleanup()
	require.NoError(t, tftest.Err)
	InPlan(tftest.Plan).That("local_file.test").Key("not_exists").DoesNotExist().IfNotFail(t)
}

func TestInSubdir(t *testing.T) {
	t.Parallel()

	tftest := setuptest.Dirs("testdata/test-in-subdir", "subdir").WithVars(nil).InitPlanShow(t)
	require.NoError(t, tftest.Err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("module.test.local_file.test").Key("content").HasValue("test").IfNotFail(t)
}
