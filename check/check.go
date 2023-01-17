package check

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

type ThatType struct {
	ResourceName string
}

func That(resourceName string) ThatType {
	return ThatType{
		ResourceName: resourceName,
	}
}

func (t ThatType) Key(key string) ThatTypeWithKey {
	return ThatTypeWithKey{
		ResourceName: t.ResourceName,
		Key:          key,
	}
}

type ThatTypeWithKey struct {
	ResourceName string
	Key          string
}

func (twk ThatTypeWithKey) HasValue(value string) ThatTypeWithKeyAndValue {
	return ThatTypeWithKeyAndValue{
		ResourceName: twk.ResourceName,
		Key:          twk.Key,
		Value:        value,
	}

}

type ThatTypeWithKeyAndValue struct {
	ResourceName string
	Key          string
	Value        string
}

func (twkav ThatTypeWithKeyAndValue) InPlan(plan terraform.PlanStruct) error {
	if _, ok := plan.ResourcePlannedValuesMap[twkav.ResourceName]; !ok {
		return fmt.Errorf(
			"%s: resource not found in plan",
			twkav.ResourceName,
		)
	}
	resource := plan.ResourcePlannedValuesMap[twkav.ResourceName]
	if _, ok := resource.AttributeValues[twkav.Key]; !ok {
		return fmt.Errorf(
			"%s: key %s not found in resource",
			twkav.ResourceName,
			twkav.Key,
		)
	}
	attrValue := resource.AttributeValues[twkav.Key]
	if assert.ObjectsAreEqual(resource.AttributeValues[twkav.Key], twkav.Value) {
		return fmt.Errorf(
			"%s: attribute %s, planned value %s not equal to assertion %s",
			twkav.ResourceName,
			twkav.Key,
			attrValue,
			twkav.Value)
	}
	return nil
}
