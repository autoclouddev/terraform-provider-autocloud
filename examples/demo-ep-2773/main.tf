terraform {
  required_providers {
    autocloud = {
      version = "0.2.0"
      source  = "autocloud.io/autoclouddev/autocloud"
    }
  }
}

provider "autocloud" {
  # endpoint = "https://api.autocloud.domain.com/api/v.0.0.1"
}

data "autocloud_github_repos" "repos" {}

####
# Local variables
locals {
  # Destination repos where generated code will be submitted
  dest_repos = [
    for repo in data.autocloud_github_repos.repos.data[*].url : repo
    if length(regexall("/terraform-generator-test", repo)) > 0
    # if length(regexall("/infrastructure-live", repo)) > 0 || length(regexall("/self-hosted-infrastructure-live", repo)) > 0
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
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/security/kms?ref=0.3.2"

  #   display_order = ["bucket"]
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
}



# ####
# # S3 Bucket
# #
# # Stores static site content
resource "autocloud_module" "s3_bucket" {
  name   = "s3bucket"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/storage/s3/bucket?ref=0.3.2"

  #   display_order = ["bucket"]
}

data "autocloud_blueprint_config" "s3_bucket_processor" {
  source = {
    s3 = autocloud_module.s3_bucket.blueprint_config
  }

  omit_variables = [
    "additional_policy_template_filename",
    "authorized_users",
    "block_public",
    "enable_bucket_ownership_controls",
    "error_file",
    "external_policy",
    "force_destroy",
    "index_file",
    "is_web_app",
    "kms_key_arn",
    "object_ownership",
    "policy_statements"
  ]

  ###
  # Force bucket owner controls
  variable {
    name  = "enable_bucket_ownership_controls"
    value = true
  }
}




# ####
# # CloudFront Distribution
# #
# # Serves the content to the public internet and configures access to the private S3 bucket
resource "autocloud_module" "cloudfront" {
  name   = "cloudfront"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/network/cloudfront/distro?ref=0.4.0"

  #   display_order = ["bucket"]
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

  ###
  ## NOTE: currently we only support conditionals targeting radio inputs
  variable {
    name = "environment"
    type = "radio"
  }
  ###


  # ###
  # # Set default object
  variable {
    name  = "default_root_object"
    value = "index.html"
  }

  ###
  # Set comment
  variable {
    name  = "comment"
    value = format("Distro for %s.%s.unstyl.com", "static-generator", "nonprod") # Need to refer to the name value here as well
  }

  ###
  # Set S3 source information
  variable {
    name  = "s3_bucket_name"
    value = autocloud_module.s3_bucket.outputs["bucket_name"]
  }

  variable {
    name  = "s3_bucket_domain_name"
    value = autocloud_module.s3_bucket.outputs["bucket_domain_name"]
  }

  # ###
  # # Set SSL Certificate based on environment
  variable {
    name         = "ssl_certificate_arn"
    display_name = "SSL Certificate ARN"
    helper_text  = "ACM Cert ARN for the environtment targeted2"

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

  # ###
  # # Set domain name
  variable {
    name = "alternate_domain_names"
    value = jsonencode([
      format("%s.%s.unstyl.com", "static-generator", "nonprod") # Need to refer to the name value here as well
    ])
  }
}



# ####
# # S3 Bucket Policy
# #
# # Enforces encryption in transit and at rest, forces KMS encryption, and allows access to the CloudFront Origin Identity
resource "autocloud_module" "s3_bucket_policy" {
  name   = "s3bucketpolicy"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/storage/s3/bucket_policy?ref=0.3.2"

  #   display_order = ["bucket"]
}

data "autocloud_blueprint_config" "s3_bucket_policy_processor" {
  source = {
    # s3_bucket = autocloud_module.s3_bucket.blueprint_config
    s3_bucket_policy = autocloud_module.s3_bucket_policy.blueprint_config
  }

  omit_variables = [
    "policies",
    "s3_bucket_name"
  ]

  ###
  # Set S3 Bucket Name
  variable {
    name  = "s3_bucket_name"
    value = autocloud_module.s3_bucket.outputs["bucket_name"]
  }

  ###
  # Set S3 Bucket Policies
  variable {
    name = "policies"
    value = jsonencode([
      autocloud_module.s3_bucket.outputs["policy"],
      autocloud_module.cloudfront.outputs["s3_bucket_policy"]
    ])
  }
}



# ####
# # Route53 Record
# #
# # Configures an ALIAS record in Route53 to serve the content on a custom URL
resource "autocloud_module" "route_53_record" {
  name   = "route53record"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/network/dns/record?ref=0.5.0"

  #   display_order = ["bucket"]
}

data "autocloud_blueprint_config" "route_53_record_processor" {
  source = {
    # s3_bucket = autocloud_module.s3_bucket.blueprint_config
    s3_bucket_policy = autocloud_module.route_53_record.blueprint_config
  }

  omit_variables = [
    "allow_overwrite",
    "s3_bucket_name"
  ]


  # Need to override these variables:
  #
  # "hostname"
  # "type"
  # "records"
  # "ttl"
  # "zone_id"
  #
  # With these values:
  #
  # hostname = format("%s.%s.unstyl.com", "static-test", "nonprod")
  # type     = "A"
  # alias = {
  #   name                   = module.unstyl_nonprod_static_test_cloudfront_distro.domain_name
  #   zone_id                = module.unstyl_nonprod_static_test_cloudfront_distro.hosted_zone_id
  #   evaluate_target_health = true
  # }
  # records = []
  # ttl     = 300
  # zone_id = "Z04736112P360P8GZTZCT"
}



####
# Create Blueprint Config
#
# Combine resources into the final config
data "autocloud_blueprint_config" "final" {
  source = {
    kms        = data.autocloud_blueprint_config.kms_key_processor.blueprint_config,
    s3         = data.autocloud_blueprint_config.s3_bucket_processor.blueprint_config
    cloudfront = data.autocloud_blueprint_config.cloudfront_processor.blueprint_config
    s3_policy  = data.autocloud_blueprint_config.s3_bucket_policy_processor.blueprint_config
    route53    = data.autocloud_blueprint_config.route_53_record_processor.blueprint_config
  }

  omit_variables = [
    "enabled"
  ]

  variable {
    name  = "enabled"
    value = true
  }

  variable {
    name         = "namespace"
    display_name = "Namespace"
    helper_text  = "The organization namespace the assets will be deployed in"

    value = "unstyl"
  }

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
    validation_rule {
      rule          = "isRequired"
      error_message = "invalid"
    }
  }

  variable {
    name         = "name"
    display_name = "Name"
    helper_text  = "The name of the static site"

    type = "shortText"

    validation_rule {
      rule          = "isRequired"
      error_message = "invalid"
    }
  }

  # variable {
  #   name    = "tags"
  #   display_name = "Tags"
  #   helper_text  = "A map of tags to apply to the deployed assets"

  #   type = "shortText"

  #   validation_rule {
  #     rule          = "isRequired"
  #     error_message = "invalid"
  #   }
  # }
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

  step 1: step-1-description
  step 2: step-2-description
  step 3: step-3-description
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
      title                   = "[AutoCloud] new Static Site, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new Static Site, created by {{authorName}}"
      body                    = file("../files/pull_request.md.tpl")
      variables = {
        authorName = "generic.authorName"
        # siteName  = "" # Need to populate this with values from global form config
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
      # namespace    = "kmskey.namespace"
      # environment = "kmskey.environment"
      # name        = "kmskey.name"
      namespace   = "unstyl"
      environment = "nonprod"
      name        = "static-generator"
    }

    modules = [
      autocloud_module.kms_key.name,
      autocloud_module.s3_bucket.name,
      autocloud_module.cloudfront.name,
      autocloud_module.s3_bucket_policy.name,
      autocloud_module.route_53_record.name
    ]
  }
}
