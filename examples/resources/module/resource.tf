terraform {
  required_providers {
    autocloud = {
      version = "0.2.0"
      source  = "autocloud.io/autoclouddev/autocloud"
    }
  }
}

resource "autocloud_module" "kms" {
  name   = "kmsRaw"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/security/kms?ref=0.1.0"
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