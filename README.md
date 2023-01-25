# terratest-terraform-fluent

Terratest extension package for testing Terraform code with fluent assertions.

## Usage

```go
package test

import (
  "testing"

  "github.com/matt-FFFFFF/terratest-terraform-fluent/check"
  "github.com/matt-FFFFFF/terratest-terraform-fluent/setuptest"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
)

const(
  basicTestData = "testdata/basic"
)

func TestSomeTerraform( t *testing.T) {
  // Set up the Terraform test and run terraform init, plan and show,
  // saving the plan output to a struct.
  // The returned struct in tftest contains the plan struct, the clean up func and any errors.
  //
  // The Dirs inputs are the test root directory and the relative path to the test code.
  // (this must be a subdirectory of the test root directory)
  // The WithVars inputs are the Terraform variables to pass to the test.
  // The InitAndPlanAndShowWithStruct inputs are the testint.T value
  tftest := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
  require.NoError(t, tftest.Err)
  defer tftest.Cleanup()

  // Check that the plan contains the expected number of resources.
  check.InPlan(plan).NumberOfResourcesEquals(1).IfNotFail(t)

  // Check that the plan contains the expected resource, with an attribute called `my_attribute` and
  // a corresponding value of `my_value`.
  check.InPlan(plan).That("my_terraform_resource.name").Key("my_attribute").HasValue("my_value").IfNotFail(t)
}
```
