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
# Local vars
locals {
  # Destination repos where generated code will be submitted
  dest_repos = [
    for repo in data.autocloud_github_repos.repos.data[*].url : repo
    if length(regexall("/terraform-generator-test", repo)) > 0
    # if length(regexall("/terraform-generator-test", repo)) > 0 || length(regexall("/infrastructure-live-demo", repo)) > 0 || length(regexall("/self-hosted-infrastructure-live", repo)) > 0
  ]
}

resource "autocloud_module" "s3_bucket" {
  name   = "s3bucket"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/storage/s3/bucket?ref=0.10.2"
}


data "autocloud_blueprint_config" "s3_custom_form" {
  source = {
    s3 = autocloud_module.s3_bucket.blueprint_config
  }

  omit_variables = [
    "acl",
    "acceleration_status",
    "bucket_prefix",
    "logging",
    "object_ownership",
    "restrict_public_buckets",
    "versioning",
    "attach_deny_insecure_transport_policy",
    "attach_elb_log_delivery_policy",
    "attach_lb_log_delivery_policy",
    "attach_policy",
    "attach_public_policy",
    "attach_require_latest_tls_policy",
    "block_public_acls",
    "block_public_policy",
    "control_object_ownership",
    "create_bucket",
    "expected_bucket_owner",
    "force_destroy",
    "ignore_public_acls",
    "object_lock_enabled",
    "policy",
    "putin_khuylo",
    "request_payer",
    "tags"
  ]
}

resource "autocloud_blueprint" "example" {
  name = "File Content Generator"

  ###
  # UI Configuration
  #
  author       = "enrique.enciso@autocloud.dev"
  description  = "Terraform Generator storage in cloud"
  instructions = "Instructions"

  labels = ["aws"]

  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new S3 maps generator, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new S3 maps generator, created by {{authorName}}"
      body                    = file("../files/pull_request.md.tpl")
      variables = {
        authorName = "generic.authorName"
      }
    }
  }


  ###
  # File definitions
  #

  file {
    action      = "CREATE"
    destination = "s3.tf"
    variables = {
    }

    modules = [
      autocloud_module.s3_bucket.name
    ]

    content = <<-EOT
        # begin - paste to the top of the generated terraform code
        terraform {
          required_version = "~> 1.1.0"
          required_providers {
            aws = {
              source  = "hashicorp/aws"
              version = "~> 4.0"
            }
          }
        }

        variable "account_num" {
          type        = string
          description = "Target AWS account number, mandatory"
        }

        variable "aws_region" {
          description = "AWS region"
          type        = string
        }

        variable "aws_role" {
          description = "AWS role to assume"
          type        = string
        }

        provider "aws" {
          region = var.aws_region
          # The following code is for using cross account assume role
          assume_role {
            role_arn = "arn:aws:iam::${var.account_num}:role/${var.aws_role}"
          }
        }

        # end - paste
        EOT

  }

  file {
    action      = "CREATE"
    destination = "s3_with_header_and_footer.tf"
    variables = {
    }

    header = "# THE HEADER IS THIS"

    footer = <<-EOT
      # begin - paste to the top of the generated terraform code
      terraform {
        required_version = "~> 1.1.0"
        required_providers {
          aws = {
            source  = "hashicorp/aws"
            version = "~> 4.0"
          }
        }
      }

      variable "account_num" {
        type        = string
        description = "Target AWS account number, mandatory"
      }

      variable "aws_region" {
        description = "AWS region"
        type        = string
      }

      variable "aws_role" {
        description = "AWS role to assume"
        type        = string
      }

      provider "aws" {
        region = var.aws_region
        # The following code is for using cross account assume role
        assume_role {
          role_arn = "arn:aws:iam::${var.account_num}:role/${var.aws_role}"
        }
      }

      # end - paste
      EOT

    modules = [
      autocloud_module.s3_bucket.name
    ]
  }

  config = data.autocloud_blueprint_config.s3_custom_form.config
}
