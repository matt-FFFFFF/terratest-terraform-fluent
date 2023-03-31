terraform {
  required_version = ">= 1.4.0"
}

locals {
  now = timestamp()
}

resource "terraform_data" "test" {
  input = local.now
}
