data "autocloud_github_repos" "repos" {}

####
# Local variables
locals {
  # Destination repos where generated code will be submitted
  dest_repos = [
    for repo in data.autocloud_github_repos.repos.data[*].url : repo
    if length(regexall("/infrastructure-live-demo", repo)) > 0
  ]
}



####
# Module Resources
#
# Connect to the Terraform modules that will be used to create this blueprint


####
# KMS Key
#
resource "autocloud_module" "kms_key" {
  name    = "cpkmskey"
  source  = "cloudposse/kms-key/aws"
  version = "0.12.1"
}

data "autocloud_blueprint_config" "kms_key_processor" {
  source = {
    kms = autocloud_module.kms_key.blueprint_config
  }

  ###
  # Hide variables from user
  omit_variables = [
    # Global
    "context",
    "tenant",
    "stage",
    "delimiter",
    "attributes",
    "labels_as_tags",
    "additional_tag_map",
    "label_order",
    "regex_replace_chars",
    "id_length_limit",
    "label_key_case",
    "label_value_case",
    "descriptor_formats",
    "enabled",
    # Use defaults in the module (don't collect)
    "alias",
    "customer_master_key_spec",
    "key_usage",
    "multi_region",
    "policy",
    # Hard Coded values
    "deletion_window_in_days",
    "description",
    "enable_key_rotation",
  ]

  ###
  # Force KMS key deletion window to 14 days
  variable {
    name = "kms.variables.deletion_window_in_days"
    type = "shortText"

    value = 14
  }

  ###
  # Set description
  variable {
    name  = "kms.variables.description"
    value = format("KMS key for encryption of KMS encrypted S3 bucket")
  }
}



####
# S3 Bucket
#
resource "autocloud_module" "s3_bucket" {
  name    = "cps3bucket"
  source  = "cloudposse/s3-bucket/aws"
  version = "3.1.0"
}

data "autocloud_blueprint_config" "s3_bucket_processor" {
  source = {
    s3 = autocloud_module.s3_bucket.blueprint_config
  }

  ###
  # Hide variables from user
  omit_variables = [
    # Global
    "context",
    "tenant",
    "stage",
    "delimiter",
    "attributes",
    "labels_as_tags",
    "additional_tag_map",
    "label_order",
    "regex_replace_chars",
    "id_length_limit",
    "label_key_case",
    "label_value_case",
    "descriptor_formats",
    "enabled",
    # Use defaults in the module (don't collect)
    "access_key_enabled",
    "acl",
    "allowed_bucket_actions",
    "block_public_acls",
    "block_public_policy",
    "bucket_key_enabled",
    "bucket_name",
    "cors_configuration",
    "force_destroy",
    "grants",
    "ignore_public_acls",
    "lifecycle_configuration_rules",
    "lifecycle_rule_ids",
    "lifecycle_rules",
    "logging",
    "object_lock_configuration",
    "policy",
    "privileged_principal_actions",
    "privileged_principal_arns",
    "replication_rules",
    "restrict_public_buckets",
    "s3_replica_bucket_arn",
    "s3_replication_enabled",
    "s3_replication_permissions_boundary_arn",
    "s3_replication_rules",
    "s3_replication_source_roles",
    "source_policy_documents",
    "ssm_base_path",
    "store_access_key_in_ssm",
    "transfer_acceleration_enabled",
    "user_enabled",
    "user_permissions_boundary_arn",
    "versioning_enabled",
    "website_configuration",
    "website_redirect_all_requests_to",
    # Hard Coded values
    "allow_encrypted_uploads_only",
    "allow_ssl_requests_only",
    "kms_master_key_arn",
    "s3_object_ownership",
    "sse_algorithm",
  ]

  ###
  # Force encrypted uploads
  variable {
    name  = "s3.variables.allow_encrypted_uploads_only"
    value = true
  }

  ###
  # Force encrypted downloads
  variable {
    name  = "s3.variables.allow_ssl_requests_only"
    value = true
  }

  ###
  # Force BucketOwner object permissions
  variable {
    name  = "s3.variables.s3_object_ownership"
    value = "BucketOwnerEnforced"
  }

  ###
  # Use KMS key encryption
  variable {
    name  = "s3.variables.sse_algorithm"
    value = "aws:kms"
  }

  ###
  # Set KMS Key ARN
  variable {
    name  = "s3.variables.kms_master_key_arn"
    value = autocloud_module.kms_key.outputs["key_arn"]
  }
}



####
# Create Blueprint Config
#
# Combine resources into the final config
data "autocloud_blueprint_config" "global" {
  source = {
    kms = data.autocloud_blueprint_config.kms_key_processor.blueprint_config,
    s3  = data.autocloud_blueprint_config.s3_bucket_processor.blueprint_config
  }

  ###
  # Hide variables from user
  omit_variables = [
    # Global
    # Use defaults in the module (don't collect)
    "context",
    "tenant",
    "stage",
    "delimiter",
    "attributes",
    "labels_as_tags",
    "additional_tag_map",
    "label_order",
    "regex_replace_chars",
    "id_length_limit",
    "label_key_case",
    "label_value_case",
    "descriptor_formats",
    # Hard Coded values
    "enabled",

    # KMS Key
    # Use defaults in the module (don't collect)
    "alias",
    "customer_master_key_spec",
    "key_usage",
    "multi_region",
    "policy",
    # Hard Coded values
    "deletion_window_in_days",
    "description",
    "enable_key_rotation",

    # S3 Bucket
    # Use defaults in the module (don't collect)
    "access_key_enabled",
    "acl",
    "allowed_bucket_actions",
    "block_public_acls",
    "block_public_policy",
    "bucket_key_enabled",
    "bucket_name",
    "cors_configuration",
    "force_destroy",
    "grants",
    "ignore_public_acls",
    "lifecycle_configuration_rules",
    "lifecycle_rule_ids",
    "lifecycle_rules",
    "logging",
    "object_lock_configuration",
    "policy",
    "privileged_principal_actions",
    "privileged_principal_arns",
    "replication_rules",
    "restrict_public_buckets",
    "s3_replica_bucket_arn",
    "s3_replication_enabled",
    "s3_replication_rules",
    "s3_replication_source_roles",
    "source_policy_documents",
    "ssm_base_path",
    "store_access_key_in_ssm",
    "transfer_acceleration_enabled",
    "user_enabled",
    "versioning_enabled",
    "website_configuration",
    "website_redirect_all_requests_to",
    # Hard Coded values
    "allow_encrypted_uploads_only",
    "allow_ssl_requests_only",
    "kms_master_key_arn",
    "s3_object_ownership",
    "sse_algorithm",
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
  # Collect the name of the asset group
  variable {
    name         = "name"
    display_name = "Name"
    helper_text  = "The name of the encrypted S3 bucket"

    type = "shortText"

    validation_rule {
      rule          = "isRequired"
      error_message = "You must provide a name for the encrypted S3 bucket"
    }
  }

  ###
  # Collect tags to apply to assets
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
  name = "KMS Encrypted S3 Bucket"

  ###
  # UI Configuration
  #
  author       = "jim@unstyl.com"
  description  = "Deploys a KMS Encrypted S3 Bucket to AWS"
  instructions = <<-EOT
  To deploy this generator, these simple steps:

    * step 1: Choose the target environment
    * step 2: Provide a name to identify assets
    * step 3: Add tags to apply to assets
  EOT

  labels = ["aws"]



  ###
  # Form configuration
  config = data.autocloud_blueprint_config.global.config



  ###
  # File definitions
  #
  file {
    action      = "CREATE"
    destination = "aws/{{environment}}/{{namespace}}-{{environment}}-{{name}}.tf"
    variables = {
      namespace   = "cpkmskey.namespace"
      environment = "cpkmskey.environment"
      name        = "cpkmskey.name"
    }

    modules = [
      autocloud_module.kms_key.name,
      autocloud_module.s3_bucket.name,
    ]
  }



  ###
  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new KMS Encrypted S3 Bucket {{namespace}}-{{environment}}-{{name}}, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new KMS Encrypted S3 Bucket {{namespace}}-{{environment}}-{{name}}, created by {{authorName}}"
      body                    = file("./files/pull_request.md.tpl")
      variables = {
        authorName  = "generic.authorName"
        namespace   = "cpkmskey.namespace"
        environment = "cpkmskey.environment"
        name        = "cpkmskey.name"
      }
    }
  }
}
