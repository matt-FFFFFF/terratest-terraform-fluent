package tfutils

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// GetDefaultTerraformOptions returns the default Terraform options for the
// given directory.
func GetDefaultTerraformOptions(t *testing.T, dir string) *terraform.Options {
	if dir == "" {
		dir = "testdata/" + t.Name()
	}
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	o := terraform.Options{
		Logger:       logger.TestingT,
		PlanFilePath: "tfplan",
		TerraformDir: dir,
		Lock:         false,

		Vars: make(map[string]interface{}),
	}
	return terraform.WithDefaultRetryableErrors(t, &o)
}
