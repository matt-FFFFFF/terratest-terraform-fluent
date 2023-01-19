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

func TestSomeTerraform( t *testing.T) {
  // Set up the Terraform options and run terraform init and plan,
  // saving the plan output to a variable.
  // The directory should be relative to the running test.
  tftest := setuptest.Dirs(basicTestData, "").WithVars(nil).InitAndPlanAndShowWithStruct(t)
  require.NoError(t, tftest.Err)
  defer tftest.Cleanup()
  // Check that the plan contains the expected number of resources.
  assert.NoError(
    t,
    check.InPlan(plan).NumberOfResourcesEquals(1)
  )

  // Check that the plan contains the expected resource, with an attribute called `my_attribute` and
  // a corresponding value of `my_value`.
  assert.NoError(
    t,
    check.InPlan(plan).That("my_terraform_resource.name").Key("my_attribute").HasValue("my_value"),
  )
}
```
