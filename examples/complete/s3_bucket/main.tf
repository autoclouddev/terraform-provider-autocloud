data "autocloud_github_repos" "repos" {
}

####
# Local variables
locals {
  # Destination repos where generated code will be submitted
  dest_repos = [
      data.autocloud_github_repos.repos.data
    ]
}



####
# Module Resources
#
# Connect to the Terraform modules that will be used to create this generator


####
# KMS Key
#
resource "autocloud_module" "kms_key" {
  name   = "kmskey"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/security/kms/key?ref=0.10.0"
}

data "autocloud_blueprint_config" "kms_key_processor" {
  source = {
    kms = autocloud_module.kms_key.blueprint_config
  }

  omit_variables = [
    # Use defaults in the module (don't collect)
    "customer_master_key_spec",
    "key_usage",
    "policies",
    "tags",

    # Hard coded
    "description",
    "enable_key_rotation",
  ]

  ###
  # Force key rotation
  variable {
    name         = "kms.variables.enable_key_rotation"
    display_name = "Automatic Key Rotation"
    helper_text  = "Enable automatic key rotation for the KMS key"

    value = true
  }

  ###
  # Set description
  variable {
    name         = "kms.variables.description"
    display_name = "KMS Key description"
    value = format("KMS key for encryption of test s3 bucket")
  }

  variable {
    name         = "kms.variables.deletion_window_in_days"
    display_name = "Key Deletion Window"
    helper_text  = "The KMS key deletion window in days"

    type = "shortText"
  }

  variable {
    name         = "kms.variables.multi_region"
    display_name = "Multi Region Key"
    helper_text  = "Whether or not the KMS key will be deployed as a multi region key"

    type = "radio"

    options {
      option {
        label   = "Single region key"
        value   = "true"
        checked = false
      }
      option {
        label   = "Multi region key"
        value   = "false"
        checked = true
      }
    }
  }
}



####
# S3 Bucket
#
# Stores
resource "autocloud_module" "s3_bucket" {
  name   = "s3bucket"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/storage/s3/bucket?ref=0.11.0"
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
  # Force public access block
  variable {
    name  = "s3.variables.block_public"
    value = true
  }

  ###
  # Force bucket owner controls
  variable {
    name  = "s3.variables.enable_bucket_ownership_controls"
    value = true
  }

  ###
  # Use an external policy
  variable {
    name  = "s3.variables.external_policy"
    value = true
  }

  ###
  # Don't delete site contents
  variable {
    name  = "s3.variables.force_destroy"
    value = false
  }
}

####
# Create Blueprint Config
#
# Combine resources into the final config
data "autocloud_blueprint_config" "final" {
  source = {
    kms              = data.autocloud_blueprint_config.kms_key_processor.blueprint_config,
    s3_bucket        = data.autocloud_blueprint_config.s3_bucket_processor.blueprint_config,
  }

  ###
  # Hide variables from user
  omit_variables = [
    # Global
    "enabled",

    # KMS Key
    "enable_key_rotation",
    "description",
    "deletion_window_in_days",
    "multi_region",

    # S3 Bucket
    "block_public",
    "enable_bucket_ownership_controls",
    "external_policy",
    "force_destroy",

    # S3 Bucket Policy
    "s3_bucket_name",
    "policies",
  ]

  ###
  # Hard code `enabled` to true to create all assets
  variable {
    name  = "enabled"
    value = true
  }

  ###
  # Set the namespace
  variable {
    name         = "namespace"
    display_name = "Namespace"
    helper_text  = "The organization namespace the assets will be deployed in"

    type = "shortText"

    value = "example"
  }

  ###
  # Choose the environment
  variable {
    name         = "environment"
    display_name = "Environment"
    helper_text  = "The environment the assets will be deployed in"

    type = "radio"

    options {
      option {
        label   = "Dev"
        value   = "Dev"
        checked = true
      }
      option {
        label = "Production"
        value = "production"
      }
    }
  }

  ###
  # Input the name
  variable {
    name         = "name"
    display_name = "Name"
    helper_text  = "The name of the s3 bucket"

    type = "shortText"

    validation_rule {
      rule          = "isRequired"
      error_message = "This field is required"
    }
  }

  variable {
    name         = "tags"
    display_name = "Tags"
    helper_text  = "A map of tags to apply to the deployed assets"

    type = "map"
  }
}



####
# Create Blueprint
#
# Create generator blueprint that contains all the elements
resource "autocloud_blueprint" "this" {
  name = "S3 Bucket"

  ###
  # UI Configuration
  #
  author       = "example@example.com"
  description  = "Secure S3 Bucket generator"
  instructions = <<-EOT
  To deploy this generator, follow these simple steps:

    * step 1: step-1-description
    * step 2: step-2-description
    * step 3: step-3-description
  EOT

  labels = ["aws"]

  ###
  # Form configuration
  config = data.autocloud_blueprint_config.final.config


  ###
  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = [local.dest_repos[0][0].url]
    git_url_default = local.dest_repos[0][0].url // "https://github.com/autoclouddev/infrastructure-live-example"

    pull_request {
      title                   = "[AutoCloud] new S3 bucket {{namespace}}-{{environment}}-{{name}}, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new S3 bucket {{namespace}}-{{environment}}-{{name}}, created by {{authorName}}"
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
    destination = "aws/{{environment}}/{{namespace}}-{{environment}}-{{name}}.tf"
    variables = {
      namespace   = "kmskey.namespace"
      environment = "kmskey.environment"
      name        = "kmskey.name"
    }

    modules = [
      autocloud_module.kms_key.name,
      autocloud_module.s3_bucket.name,
    ]
  }
}
