package check

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// ThatType is a type which can be used for more fluent assertions for a given Resource
type ThatType struct {
	Plan         *terraform.PlanStruct
	ResourceName string
}

type JsonAssertionFunc func(input json.RawMessage) (*bool, error)

// That returns a type which can be used for more fluent assertions for a given Resource
func (p PlanType) That(resourceName string) ThatType {
	return ThatType{
		Plan:         p.Plan,
		ResourceName: resourceName,
	}
}

// Exists returns an error if the resource does not exist in the plan
func (t ThatType) Exists() *CheckError {
	if _, ok := t.Plan.ResourcePlannedValuesMap[t.ResourceName]; !ok {
		return newCheckErrorf(
			"%s: resource not found in plan",
			t.ResourceName,
		)
	}
	return nil
}

// DoesNotExist returns an error if the resource exists in the plan
func (t ThatType) DoesNotExist() *CheckError {
	if _, exists := t.Plan.ResourcePlannedValuesMap[t.ResourceName]; exists {
		return newCheckErrorf(
			"%s: resource found in plan",
			t.ResourceName,
		)
	}
	return nil
}

// Key returns a type which can be used for more fluent assertions for a given Resource & Key combination
func (t ThatType) Key(key string) ThatTypeWithKey {
	return ThatTypeWithKey{
		Plan:         t.Plan,
		ResourceName: t.ResourceName,
		Key:          key,
	}
}

// ThatTypeWithKey is a type which can be used for more fluent assertions for a given Resource & Key combination
type ThatTypeWithKey struct {
	Plan         *terraform.PlanStruct
	ResourceName string
	Key          string
}

// HasValue returns a CheckError if the resource does not exist in the plan or if the value of the key does not match the
// expected value
func (twk ThatTypeWithKey) HasValue(expected interface{}) *CheckError {
	if err := twk.Exists(); err != nil {
		return err
	}

	resource := twk.Plan.ResourcePlannedValuesMap[twk.ResourceName]
	actual := resource.AttributeValues[twk.Key]

	if err := validateEqualArgs(expected, actual); err != nil {
		return newCheckErrorf("invalid operation: %#v == %#v (%s)",
			expected,
			actual,
			err,
		)
	}

	if !assert.ObjectsAreEqualValues(actual, expected) {
		return newCheckErrorf(
			"%s: attribute %s, planned value %s not equal to assertion %s",
			twk.ResourceName,
			twk.Key,
			actual,
			expected,
		)
	}
	return nil
}

// ContainsJsonValue returns a *CheckError which asserts upon a given JSON string set into
// the State by deserializing it and then asserting on it via the JsonAssertionFunc
func (twk ThatTypeWithKey) ContainsJsonValue(assertion JsonAssertionFunc) *CheckError {
	if err := twk.Exists(); err != nil {
		return err
	}

	if twk.HasValue("") == nil {
		return newCheckErrorf(
			"%s: key %s was empty",
			twk.ResourceName,
			twk.Key,
		)
	}

	resource := twk.Plan.ResourcePlannedValuesMap[twk.ResourceName]
	actual := resource.AttributeValues[twk.Key]
	j := json.RawMessage(actual.(string))
	ok, err := assertion(j)
	if err != nil {
		return newCheckErrorf(
			"%s: asserting value for %q: %+v",
			twk.ResourceName,
			twk.Key,
			err,
		)
	}

	if ok == nil || !*ok {
		return newCheckErrorf(
			"%s: assertion failed for %q: %+v",
			twk.ResourceName,
			twk.Key,
			err,
		)
	}

	return nil
}

// Exists returns a CheckError if the resource does not exist in the plan or if the key does not exist in the resource
func (twk ThatTypeWithKey) Exists() *CheckError {
	if err := InPlan(twk.Plan).That(twk.ResourceName).Exists(); err != nil {
		return newCheckError(err.Error())
	}

	resource := twk.Plan.ResourcePlannedValuesMap[twk.ResourceName]
	if _, exists := resource.AttributeValues[twk.Key]; !exists {
		return newCheckErrorf(
			"%s: key %s not found in resource",
			twk.ResourceName,
			twk.Key,
		)
	}
	return nil
}

// DoesNotExist returns a CheckError if the resource does not exist in the plan or if the key exists in the resource
func (twk ThatTypeWithKey) DoesNotExist() *CheckError {
	if err := InPlan(twk.Plan).That(twk.ResourceName).Exists(); err != nil {
		return newCheckError(err.Error())
	}

	resource := twk.Plan.ResourcePlannedValuesMap[twk.ResourceName]
	if _, exists := resource.AttributeValues[twk.Key]; exists {
		return newCheckErrorf(
			"%s: key %s found in resource",
			twk.ResourceName,
			twk.Key,
		)
	}
	return nil
}

// validateEqualArgs checks whether provided arguments can be safely used in the
// HasValue function.
func validateEqualArgs(expected, actual interface{}) error {
	if expected == nil && actual == nil {
		return nil
	}

	if isFunction(expected) || isFunction(actual) {
		return fmt.Errorf("cannot take func type as argument")
	}
	return nil
}

func isFunction(arg interface{}) bool {
	if arg == nil {
		return false
	}
	return reflect.TypeOf(arg).Kind() == reflect.Func
}
