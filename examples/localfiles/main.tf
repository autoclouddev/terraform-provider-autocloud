terraform {
  required_providers {
    autocloud = {
      version = "0.2.0"
      source  = "autocloud.io/autoclouddev/autocloud"
    }
  }
}

####
# Module Resources
#
# Connect to the Terraform modules that will be used to create this blueprint


####
# KMS Key
#
resource "autocloud_module" "bucket" {
  name = "s3localbucket"

  ###
  # This source works as expected
  # source  = "terraform-aws-modules/kms/aws"

  ###
  # This source is a private registry on Terraform Cloud, does not work correctly
  source = "./source_module/s3/bucket"
}

data "autocloud_blueprint_config" "kms_key_processor" {
  source = {
    s3 = autocloud_module.bucket.blueprint_config
  }
}

####
# Create Blueprint
#
# Create generator blueprint that contains all the elements
resource "autocloud_blueprint" "this" {
  name = "[VALIDATION] LOCAL S3"

  ###
  # UI Configuration
  #
  author       = "jim@unstyl.com"
  description  = "Deploys a KMS Key using the AWS KMS Key Module"
  instructions = <<-EOT
  To deploy this generator, these simple steps:

    * step 1: Choose the target environment
    * step 2: Provide a name to identify assets
    * step 3: Add tags to apply to assets
  EOT

  labels = ["aws"]



  ###
  # Form configuration
  # config = autocloud_module.kms_key.blueprint_config
  # config = data.autocloud_blueprint_config.kms_key_processor.config
  config = data.autocloud_blueprint_config.kms_key_processor.config


  ###
  # File definitions
  #
  file {
    action      = "CREATE"
    destination = "aws/{{environment}}/{{namespace}}-{{environment}}-{{name}}.tf"
    variables = {
      namespace   = "namespace"
      environment = "env"
      name        = "name"
    }

    modules = [
      autocloud_module.bucket.name
    ]
  }
}




