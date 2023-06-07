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
    value = format("KMS key for encryption of test static site content")
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
# Stores static site content
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
# CloudFront Distribution
#
# Serves static site content
resource "autocloud_module" "cloudfront" {
  name   = "cloudfront"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/network/cloudfront/distro?ref=0.4.0"
}

data "autocloud_blueprint_config" "cloudfront_processor" {
  source = {
    cloudfront = autocloud_module.cloudfront.blueprint_config
  }

  omit_variables = [
    # Use defaults in the module (don't collect)
    "enable_compression",
    "http_version",
    "ipv6_enabled",
    "lambda_functions",
    "price_class",
    "ssl_policy",

    # Hard coded
    "comment",
    "default_root_object",
    "s3_bucket_name",
    "s3_bucket_domain_name",
    "ssl_certificate_arn",
  ]

  ###
  # Environment conditional hack
  variable {
    name = "environment"
    type = "radio"
  }

  ###
  # Set alternate_domain_names
  variable {
    name         = "cloudfront.variables.alternate_domain_names"
    display_name = "CloudFront allowed hostnames"
    helper_text  = "List of hostnames CloudFront will serve content from, should match the Route53 record hostname"

    type = "checkbox"

    options {

    }
  }

  ###
  # Set comment
  variable {
    name = "cloudfront.variables.comment"
    value = format("Site content for test static site content")
  }

  ###
  # Set default_root_object
  variable {
    name  = "cloudfront.variables.default_root_object"
    value = autocloud_module.s3_bucket.outputs["index_file"]
  }

  ###
  # Set s3_bucket_name
  variable {
    name  = "cloudfront.variables.s3_bucket_name"
    value = autocloud_module.s3_bucket.outputs["bucket_name"]
  }

  ###
  # Set s3_bucket_domain_name
  variable {
    name  = "cloudfront.variables.s3_bucket_domain_name"
    value = autocloud_module.s3_bucket.outputs["bucket_domain_name"]
  }

  ###
  # Set ssl_certificate_arn
  variable {
    name = "cloudfront.variables.ssl_certificate_arn"

    type = "radio"

    conditional {
      source    = "cloudfront.environment"
      condition = "nonprod"
      content {
        value = "arn:aws:acm:us-east-1:534614196230:certificate/9b1e8d89-2f21-41b8-942b-db634f83b083"
      }
    }

    conditional {
      source    = "cloudfront.environment"
      condition = "production"
      content {
        value = ""
      }
    }
  }
}



####
# S3 Bucket Policy
#
# Sets the S3 bucket policy to allow CloudFront traffic
resource "autocloud_module" "s3_bucket_policy" {
  name   = "s3bucketpolicy"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/storage/s3/bucket_policy?ref=0.3.2"
}

data "autocloud_blueprint_config" "s3_bucket_policy_processor" {
  source = {
    policy = autocloud_module.s3_bucket_policy.blueprint_config
  }

  omit_variables = [
    # Use defaults in the module (don't collect)

    # Hard coded
    "s3_bucket_name",
    "policies"
  ]

  ###
  # Set S3 bucket name
  variable {
    name  = "policy.variables.s3_bucket_name"
    value = autocloud_module.s3_bucket.outputs["bucket_name"]
  }

  ###
  # Set S3 bucket policies
  variable {
    name = "policy.variables.policies"
    value = jsonencode([
      autocloud_module.s3_bucket.outputs["policy"],
      autocloud_module.cloudfront.outputs["s3_bucket_policy"]
    ])
  }
}



####
# Route53 Record
#
# Sets the Route53 record for the static site
resource "autocloud_module" "route53" {
  name   = "route53"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/network/dns/record?ref=0.5.0"
}

data "autocloud_blueprint_config" "route53_processor" {
  source = {
    route53 = autocloud_module.route53.blueprint_config
  }

  omit_variables = [
    # Use defaults in the module (don't collect)
    "allow_overwrite",

    # Hard coded
    "records",
    "ttl",
    "type",
    "zone_id"
  ]

  ###
  # Environment conditional hack
  variable {
    name = "environment"
    type = "radio"
  }

  ###
  # Set
  variable {
    name         = "route53.variables.hostname"
    display_name = "DNS hostname"
    helper_text  = "Route53 record hostname, should match list of hostnames CloudFront will serve content from"

    type = "shortText"
  }

  ###
  # Set records
  variable {
    name  = "route53.variables.records"
    value = jsonencode([])
  }

  ###
  # Set record Alias configuration
  variable {
    name = "route53.variables.alias"
    value = jsonencode({
      name                   = autocloud_module.cloudfront.outputs["domain_name"]
      zone_id                = autocloud_module.cloudfront.outputs["hosted_zone_id"]
      evaluate_target_health = true
    })
  }

  ###
  # Set record TTL
  variable {
    name  = "route53.variables.ttl"
    value = 3600
  }

  ###
  # Set record type
  variable {
    name  = "route53.variables.type"
    value = "A"
  }

  ###
  # Set zone ID
  variable {
    name = "route53.variables.zone_id"

    type = "radio"

    conditional {
      source    = "route53.environment"
      condition = "nonprod"
      content {
        value = "Z04736112P360P8GZTZCT"
      }
    }

    conditional {
      source    = "route53.environment"
      condition = "production"
      content {
        value = ""
      }
    }
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
    cloudfront       = data.autocloud_blueprint_config.cloudfront_processor.blueprint_config,
    s3_bucket_policy = data.autocloud_blueprint_config.s3_bucket_policy_processor.blueprint_config,
    route53          = data.autocloud_blueprint_config.route53_processor.blueprint_config,
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

    # CloudFront Distribution
    "comment",
    "default_root_object",
    "s3_bucket_name",
    "s3_bucket_domain_name",
    "ssl_certificate_arn",

    # S3 Bucket Policy
    "s3_bucket_name",
    "policies",

    # Route53 Module
    "alias",
    "allow_overwrite",
    # "hostname",
    "records",
    "ttl",
    "type",
    "zone_id"
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
    helper_text  = "The name of the static site"

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
  name = "Static Site"

  ###
  # UI Configuration
  #
  author       = "chris@autocloud.dev"
  description  = "Secure static site generator"
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
      title                   = "[AutoCloud] new Static Site {{namespace}}-{{environment}}-{{name}}, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new Static Site {{namespace}}-{{environment}}-{{name}}, created by {{authorName}}"
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
      autocloud_module.cloudfront.name,
      autocloud_module.s3_bucket_policy.name,
      autocloud_module.route53.name,
    ]
  }
}
