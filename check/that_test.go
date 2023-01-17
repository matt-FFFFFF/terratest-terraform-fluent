package check

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/matt-FFFFFF/terratest-terraform-fluent/tfutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	basicTestData = "testdata/basic"
)

func TestHasValueStrings(t *testing.T) {
	t.Parallel()

	opts := tfutils.GetDefaultTerraformOptions(t, basicTestData)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, opts)
	require.NoError(t, err)
	assert.NoError(
		t,
		InPlan(plan).That("local_file.test").Key("content").HasValue("test"),
	)
}

func TestHasValueStringsNotEqualError(t *testing.T) {
	t.Parallel()

	opts := tfutils.GetDefaultTerraformOptions(t, basicTestData)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, opts)
	require.NoError(t, err)
	assert.ErrorContains(
		t,
		InPlan(plan).That("local_file.test").Key("content").HasValue("throwError"),
		"attribute content, planned value test not equal to assertion throwError",
	)
}

func TestHasValueStringsToInt(t *testing.T) {
	t.Parallel()

	opts := tfutils.GetDefaultTerraformOptions(t, basicTestData)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, opts)
	require.NoError(t, err)
	assert.Error(
		t,
		InPlan(plan).That("local_file.test_int").Key("content").HasValue(123),
	)
}

func TestKeyNotExistsError(t *testing.T) {
	t.Parallel()

	opts := tfutils.GetDefaultTerraformOptions(t, basicTestData)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, opts)
	require.NoError(t, err)
	assert.ErrorContains(
		t,
		InPlan(plan).That("local_file.test").Key("not_exists").Exists(),
		"key not_exists not found in resource",
	)
}

func TestKeyNotExists(t *testing.T) {
	t.Parallel()

	opts := tfutils.GetDefaultTerraformOptions(t, basicTestData)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, opts)
	require.NoError(t, err)
	assert.NoError(
		t,
		InPlan(plan).That("local_file.test").Key("not_exists").DoesNotExist(),
	)
}
