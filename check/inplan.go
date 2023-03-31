package check

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/matt-FFFFFF/terratest-terraform-fluent/testerror"
)

func InPlan(plan *terraform.PlanStruct) PlanType {
	return PlanType{
		Plan: plan,
	}
}

type PlanType struct {
	Plan *terraform.PlanStruct
}

func (p PlanType) NumberOfResourcesEquals(expected int) *testerror.Error {
	actual := len(p.Plan.ResourcePlannedValuesMap)
	if actual != expected {
		return testerror.Newf("expected %d resources, got %d", expected, actual)
	}
	return nil
}
