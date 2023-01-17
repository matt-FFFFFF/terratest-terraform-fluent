package check

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/matt-FFFFFF/terratest-ext/tfutils"
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

func TestHasValueStringsNotEqual(t *testing.T) {
	t.Parallel()

	opts := tfutils.GetDefaultTerraformOptions(t, basicTestData)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, opts)
	require.NoError(t, err)
	assert.Error(
		t,
		InPlan(plan).That("local_file.test").Key("content").HasValue("test2"),
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
