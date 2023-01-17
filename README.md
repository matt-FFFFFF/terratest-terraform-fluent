# terratest-terraform-fluent

Terratest extension package for testing Terraform code with fluent assertions.

## Usage

```go
package test

import (
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"

  "github.com/matt-FFFFFF/terratest-terraform-fluent/check"
  "github.com/matt-FFFFFF/terratest-terraform-fluent/tfutils"
  "github.com/gruntwork-io/terratest/modules/terraform"
)

func TestSomeTerraform( t *testing.T) {
  opts := tfutils.GetDefaultTerraformOptions(t, "testdata/my-directory")
  plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
  require.NoError(t, err)

  assert.NoError(
    t,
    check.InPlan(plan).NumberOfResourcesEquals(1)
  )

  assert.NoError(
    t,
    check.InPlan(plan).That("my_terraform_resource.name").Key("my_attribute").HasValue("my_value"),
  )
}
```
