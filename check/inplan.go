package check

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func InPlan(plan *terraform.PlanStruct) PlanType {
	return PlanType{
		Plan: plan,
	}
}

type PlanType struct {
	Plan *terraform.PlanStruct
}

func (p PlanType) NumberOfResourcesEquals(expected int) error {
	actual := len(p.Plan.ResourcePlannedValuesMap)
	if actual != expected {
		return fmt.Errorf("expected %d resources, got %d", expected, actual)
	}
	return nil
}
