package check

import (
	"fmt"
	"reflect"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func InPlan(plan *terraform.PlanStruct) PlanType {
	return PlanType{
		Plan: plan,
	}
}

type PlanType struct {
	Plan *terraform.PlanStruct
}

type ThatType struct {
	Plan         *terraform.PlanStruct
	ResourceName string
}

func (p PlanType) That(resourceName string) ThatType {
	return ThatType{
		Plan:         p.Plan,
		ResourceName: resourceName,
	}
}

func (t ThatType) Exists() error {
	if _, ok := t.Plan.ResourcePlannedValuesMap[t.ResourceName]; !ok {
		return fmt.Errorf(
			"%s: resource not found in plan",
			t.ResourceName,
		)
	}
	return nil
}

func (t ThatType) Key(key string) ThatTypeWithKey {
	return ThatTypeWithKey{
		Plan:         t.Plan,
		ResourceName: t.ResourceName,
		Key:          key,
	}
}

type ThatTypeWithKey struct {
	Plan         *terraform.PlanStruct
	ResourceName string
	Key          string
}

func (twk ThatTypeWithKey) HasValue(expected interface{}) error {
	if err := twk.Exists(); err != nil {
		return err
	}

	resource := twk.Plan.ResourcePlannedValuesMap[twk.ResourceName]
	actual := resource.AttributeValues[twk.Key]

	if err := validateEqualArgs(expected, actual); err != nil {
		return fmt.Errorf("invalid operation: %#v == %#v (%s)",
			expected,
			actual,
			err,
		)
	}

	if !assert.ObjectsAreEqualValues(actual, expected) {
		return fmt.Errorf(
			"%s: attribute %s, planned value %s not equal to assertion %s",
			twk.ResourceName,
			twk.Key,
			actual,
			expected,
		)
	}
	return nil
}

func (twk ThatTypeWithKey) Exists() error {
	if err := InPlan(twk.Plan).That(twk.ResourceName).Exists(); err != nil {
		return err
	}

	resource := twk.Plan.ResourcePlannedValuesMap[twk.ResourceName]
	if _, exists := resource.AttributeValues[twk.Key]; !exists {
		return fmt.Errorf(
			"%s: key %s not found in resource",
			twk.ResourceName,
			twk.Key,
		)
	}

	return nil
}

func (twk ThatTypeWithKey) DoesNotExist() error {
	if err := InPlan(twk.Plan).That(twk.ResourceName).Exists(); err != nil {
		return err
	}

	resource := twk.Plan.ResourcePlannedValuesMap[twk.ResourceName]
	if _, exists := resource.AttributeValues[twk.Key]; exists {
		return fmt.Errorf(
			"%s: key %s found in resource",
			twk.ResourceName,
			twk.Key,
		)
	}

	return nil
}

// validateEqualArgs checks whether provided arguments can be safely used in the
// HasValue/NotEqual functions.
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
