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
# Module Resources
#
# Connect to the Terraform modules that will be used to create this generator


####
# KMS Key
#
resource "autocloud_module" "kms_key" {
  name   = "kmskey"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/security/kms/key?ref=0.10.2"
}

data "autocloud_blueprint_config" "kms_key_processor" {
  source = {
    kms = autocloud_module.kms_key.blueprint_config
  }

  omit_variables = [
    # Use defaults in the module (don't collect)
    "customer_master_key_spec",
    "description",
    "enable_key_rotation",
    "key_usage",
    "policies",
    "tags"
  ]

  ###
  # Force key rotation
  variable {
    name  = "kms.variables.enable_key_rotation"
    display_name = "Automatic Key Rotation"
    helper_text  = "Enable automatic key rotation for the KMS key"

    value = true
  }

  ###
  # Set description
  variable {
    name = "kms.variables.description"
    display_name = "KMS Key description"
    value = format("KMS key for encryption of KMS encrypted S3 bucket")
  }

  variable {
    name = "kms.variables.deletion_window_in_days"
    type = "shortText"

    value = 14
  }

  variable {
    name = "kms.variables.multi_region"
    type = "shortText"

    value = false
  }
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
    "index_file",
    "is_web_app",
    "object_ownership",
    "policy_statements",

    # Hard coded
    "block_public",
    "enable_bucket_ownership_controls",
    "external_policy",
    "force_destroy",
    "lifecycle_rules",
    
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
  # Link KMS Key arn
  variable {
    name  = "s3.variables.kms_key_arn"
    value = autocloud_module.kms_key.outputs["key_arn"]
  }

  ###
  # Don't delete site contents
  variable {
    name  = "s3.variables.force_destroy"
    value = false
  }

  ###
  # TODO: Implement Lifecycle rules
  # 
  # Output to module should be:
  #
  # lifecycle_rules = [
  #   {
  #     id      = "Retain for 1 Year"
  #     enabled = true
  #
  #     current_version_expiration = {
  #       days = 365
  #     }
  #
  #     current_version_transition = [
  #       {
  #         days          = 90
  #         storage_class = "STANDARD_IA"
  #       },
  #       {
  #         days          = 180
  #         storage_class = "GLACIER"
  #       }
  #     ]
  #
  #     noncurrent_version_expiration = {
  #       newer_noncurrent_versions = 3
  #       noncurrent_days           = 365
  #     }
  #
  #     noncurrent_version_transition = [
  #       {
  #         newer_noncurrent_versions = 3
  #         noncurrent_days           = 90
  #         storage_class             = "STANDARD_IA"
  #       },
  #       {
  #         newer_noncurrent_versions = 3
  #         noncurrent_days           = 180
  #         storage_class             = "GLACIER"
  #       }
  #     ]
  #   },
  #   {
  #     id      = "Remove expired multi-part uploads"
  #     enabled = true
  #
  #     abort_incomplete_multipart_upload_days = 7
  #   }
  # ]
}  



####
# Create Blueprint Config
#
# Combine resources into the final config
data "autocloud_blueprint_config" "final" {
  source = {
    kms = data.autocloud_blueprint_config.kms_key_processor.blueprint_config,
    s3_bucket = data.autocloud_blueprint_config.s3_bucket_processor.blueprint_config,
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

    value = "unstyl"
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
        label   = "Nonprod"
        value   = "nonprod"
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
    helper_text  = "The name of the KMS encrypted S3 bucket"

    type = "shortText"

    validation_rule {
      rule          = "isRequired"
      error_message = "This field is required"
    }
  }

  variable {
    name    = "tags"
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
  name = "KMS Encrypted S3 Bucket"

  ###
  # UI Configuration
  #
  author       = "chris@autocloud.dev"
  description  = "KMS Encrypted S3 Bucket"
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
  # TODO:
  # - Reference site name using global values
  git_config {
    destination_branch = "main"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new KMS Encrypted S3 Bucket {{namespace}}-{{environment}}-{{name}}, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new KMS Encrypted S3 Bucket {{namespace}}-{{environment}}-{{name}}, created by {{authorName}}"
      body                    = file("./files/pull_request.md.tpl")
      variables = {
        authorName = "generic.authorName"
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
