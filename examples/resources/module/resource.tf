terraform {
  required_providers {
    autocloud = {
      source  = "autoclouddev/autocloud"
      version = "~>0.9"
    }
  }
}

resource "autocloud_module" "kms" {
  name    = "cpkmskey"
  source  = "cloudposse/kms-key/aws"
  version = "0.12.1"

  header = <<-EOT
  providers = {
    aws = aws.usw3
  }
  EOT

  footer = <<-EOT
  depends_on = [
    module.account_baseline # Force account baseline before creating keys
  ]
  EOT
}
