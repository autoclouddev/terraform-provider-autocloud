terraform {
  required_providers {
    autocloud = {
      version = "0.2.0"
      source  = "autocloud.io/autoclouddev/autocloud"
    }
  }
}

data "autocloud_github_repos" "repos" {}

####
# Local variables
locals {
  # Destination repos where generated code will be submitted
  dest_repos = [
    for repo in data.autocloud_github_repos.repos.data[*].url : repo
    if length(regexall("/infrastructure-live-demo", repo)) > 0 || length(regexall("/self-hosted-infrastructure-live", repo)) > 0
  ]
}

####
# S3 Bucket
#
# Stores static site content
resource "autocloud_module" "s3_bucket" {
  name   = "s3bucket"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/storage/s3/bucket?ref=0.10.2"
}

data "autocloud_blueprint_config" "s3_bucket_processor" {
  source = {
    s3 = autocloud_module.s3_bucket.blueprint_config
  }

  omit_variables = [
    # Use defaults in the module (don't collect)
    "additional_policy_template_filename",
    "additional_policy_vars",
    "authorized_users",
    "cors_rules",
    "error_file",
    "is_web_app",
    "object_ownership",
    "policy_statements",
    "lifecycle_rules",
    # Hard coded
    "block_public",
    "enable_bucket_ownership_controls",
    "external_policy",
    "force_destroy",

    # To be configured later
    "kms_key_arn",
  ]

  ###
  # Don't delete site contents
  variable {
    name  = "s3.variables.force_destroy"
    value = "false"
  }

  display_order  {
    priority = 0
    values = ["name", "s3.variables.environment", "s3.variables.force_destroy"]
  }
}

####
# KMS Key
#
resource "autocloud_module" "kms_key" {
  name   = "kmskey"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/security/kms?ref=0.3.2"
}

data "autocloud_blueprint_config" "kms_key_processor" {
  source = {
    kms = autocloud_module.kms_key.blueprint_config
  }

  ###
  # Force key rotation
  variable {
    name  = "enable_key_rotation"
    value = true
  }

  ###
  # Set description
  variable {
    name = "description"
    # value   = "KMS key for encryption of {{namespace}}-{{environment}}-{{name}}"
    value = format("KMS key for encryption of test static site content")
  }

  display_order  {
    priority = 1
    values = ["name", "kms.variables.description"]
  }
}

# ####
# # CloudFront Distribution
# #
# # Serves the content to the public internet and configures access to the private S3 bucket
resource "autocloud_module" "cloudfront" {
  name   = "cloudfront"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/network/cloudfront/distro?ref=0.4.0"
}

data "autocloud_blueprint_config" "cloudfront_processor" {
  source = {
    cloudfront = autocloud_module.cloudfront.blueprint_config
  }

  omit_variables = [
    "alternate_domain_names",
    "comment",
    "default_root_object",
    "enable_compression",
    "http_version",
    "ipv6_enabled",
    "price_class",
    "s3_bucket_name",
    "s3_bucket_domain_name",
    "ssl_certificate_arn",
    "ssl_policy"
  ]

  variable {
    name  = "default_root_object"
    value = "index.html"
  }

 variable {
    name  = "comment"
    value = format("Distro for %s.%s.unstyl.com", "static-generator", "nonprod") # Need to refer to the name value here as well
  }

  variable {
    name  = "s3_bucket_name"
    value = autocloud_module.s3_bucket.outputs["bucket_name"]
  }

  variable {
    name  = "s3_bucket_domain_name"
    value = autocloud_module.s3_bucket.outputs["bucket_domain_name"]
  }

  display_order {
    priority = 2
    values = [ "cloudfront.variables.comment", "s3_bucket_name", "cloudfront.variables.s3_bucket_domain_name" ]
  }
}

data "autocloud_blueprint_config" "prefinal" {
  source = {
    kms        = data.autocloud_blueprint_config.kms_key_processor.blueprint_config,
    s3         = data.autocloud_blueprint_config.s3_bucket_processor.blueprint_config
    cloudfront = data.autocloud_blueprint_config.cloudfront_processor.blueprint_config
   }

   display_order {
    priority = 3
    values = ["s3.variables.cors_rules", "cloudfront.variables.s3_bucket_domain_name", "kms.variables.enable_key_rotation"]
  }
}

####
# Create Blueprint Config
#
# Combine resources into the final config
data "autocloud_blueprint_config" "final" {
  source = {
    prefinal        = data.autocloud_blueprint_config.prefinal.blueprint_config
  }

  variable {
    name = "environment"
    type = "radio"
    options {
        option {
          label   = "dev"
          value   = "dev"
          checked = false
        }
        option {
          label   = "nonprod"
          value   = "nonprod"
          checked = true
        }
        option {
          label = "prod"
          value = "prod"
          checked = false
        }
      }
  }

  display_order {
    priority = 4
    values = ["prefinal.variables.tags"]
  }
}

####
# Create Blueprint
#
# Create generator blueprint that contains all the elements
resource "autocloud_blueprint" "this" {
  name = "Display order"
  ###
  # UI Configuration
  #
  author       = "drosas@autocloud.dev"
  description  = "Display order generator"
  instructions = <<-EOT
  To deploy this generator, follow these simple steps:
  step 1: step-1-description
  step 2: step-2-description
  step 3: step-3-description
  EOT
  labels       = ["aws"]
  ###
  # Form configuration
  config = data.autocloud_blueprint_config.final.config
  ###
  # Destination repository git configuraiton
  #
  # TODO:
  # - Reference site name using global values
  git_config {
    destination_branch = "main"
    git_url_options    = local.dest_repos
    git_url_default    = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default
    pull_request {
      title                   = "[AutoCloud] new Static Site {{namespace}}-{{environment}}-{{name}}, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new Static Site {{namespace}}-{{environment}}-{{name}}, created by {{authorName}}"
      body                    = file("../files/pull_request.md.tpl")
      variables = {
        authorName  = "generic.authorName"
        namespace   = "kmskey.namespace"
        environment = "kmskey.environment"
        name        = "kmskey.name"
      }
    }
  }
  ###
  # File definitions
  #
  file {
    action      = "CREATE"
    destination = "display_order.tf"
    variables = {}
    modules = [
      autocloud_module.s3_bucket.name,
      autocloud_module.kms_key.name,
      autocloud_module.cloudfront.name,
    ]
  }
}

## Expected display order:
  #  [
  #   "name",
  #   "s3.variables.environment",
  #   "s3.variables.force_destroy",
  #   "kms.variables.description",
  #   "cloudfront.variables.comment",
  #   "s3_bucket_name",
  #   "cloudfront.variables.s3_bucket_domain_name",
  #   "s3.variables.cors_rules",
  #   "cloudfront.variables.s3_bucket_domain_name",
  #   "kms.variables.enable_key_rotation"
  #   "prefinal.variables.tags"
  #  ]
