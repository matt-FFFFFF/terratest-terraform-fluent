package check

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/matt-FFFFFF/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	basicTestData = "testdata/basic"
)

func Bool(b bool) *bool {
	var b2 = b
	return &b2
}

func TestHasValueStrings(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test").Key("content").HasValue("test").ErrorIsNil(t)
}

func TestHasValueStringsNotEqualError(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	assert.ErrorContains(
		t,
		InPlan(tftest.Plan).That("local_file.test").Key("content").HasValue("throwError"),
		"attribute content, planned value test not equal to assertion throwError",
	)
}

func TestHasValueStringsToInt(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	assert.Error(
		t,
		InPlan(tftest.Plan).That("local_file.test_int").Key("content").HasValue(123).AsError(),
	)
}

func TestKeyNotExistsError(t *testing.T) {
	t.Parallel()

	tftest, _ := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	defer tftest.Cleanup()
	assert.ErrorContains(
		t,
		InPlan(tftest.Plan).That("local_file.test").Key("not_exists").Exists(),
		"key not_exists not found in resource",
	)
}

func TestKeyNotExists(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	defer tftest.Cleanup()
	require.NoError(t, err)
	InPlan(tftest.Plan).That("local_file.test").Key("not_exists").DoesNotExist().ErrorIsNil(t)
}

func TestInSubdir(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs("testdata/test-in-subdir", "subdir").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("module.test.local_file.test").Key("content").HasValue("test").ErrorIsNil(t)
}

func TestJsonArrayAssertionFunc(t *testing.T) {
	t.Parallel()

	f := func(input json.RawMessage) (*bool, error) {
		i := make([]interface{}, 0, 1)
		if err := json.Unmarshal(input, &i); err != nil {
			return nil, fmt.Errorf("JSON input is not an array")
		}
		if len(i) == 0 {
			return nil, fmt.Errorf("JSON input is empty")
		}
		if i[0].(map[string]interface{})["test"] != "test" {
			return nil, fmt.Errorf("JSON input key name is not equal to test")
		}

		return Bool(true), nil
	}

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test_array_json").Key("content").ContainsJsonValue(JsonAssertionFunc(f)).ErrorIsNil(t)
}

func TestJsonSimpleAssertionFunc(t *testing.T) {
	t.Parallel()

	f := JsonAssertionFunc(
		func(input json.RawMessage) (*bool, error) {
			i := make(map[string]interface{})
			if err := json.Unmarshal(input, &i); err != nil {
				return nil, fmt.Errorf("JSON input is not an map")
			}
			if len(i) == 0 {
				return nil, fmt.Errorf("JSON input is empty")
			}
			if i["test"] != "test" {
				return Bool(false), nil
			}
			return Bool(true), nil
		},
	)

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test_simple_json").Key("content").ContainsJsonValue(f).ErrorIsNil(t)
}
