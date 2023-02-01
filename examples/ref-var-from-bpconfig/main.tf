

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
  name = "S3BucketMaps"

  version = "3.4.0"
  source  = "terraform-aws-modules/s3-bucket/aws" // https://registry.terraform.io/modules/terraform-aws-modules/s3-bucket/aws/latest
}

resource "autocloud_module" "kms" {
  name   = "kms"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/security/kms?ref=0.1.0"
}

data "autocloud_blueprint_config" "kms_custom_form" {
  source = {
    kms = autocloud_module.kms.blueprint_config
  }


  variable {
    name         = "environment"
    display_name = "environment target"
    helper_text  = "environment target description"
    type         = "radio"
    options {
      option {
        label   = "dev"
        value   = "dev"
        checked = true
      }
      option {
        label = "nonprod"
        value = "nonprod"
      }
      option {
        label = "prod"
        value = "prod"
      }
    }
    validation_rule {
      rule          = "isRequired"
      error_message = "invalid"
    }
  }

  variable {
    name  = "name"
    value = "this-is-the-name"
  }

  variable {
    name  = "namespace"
    value = "this-is-the-namespace"
  }
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
    // "block_public_policy",
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

  
  variable {
    name         = "block_public_policy"
    // reference to kms_custom_form blueprint config
    value        = data.autocloud_blueprint_config.kms_custom_form.variables["enable_key_rotation"] 
  }

  variable {
    name  = "s3_bucket_name"
    value = data.autocloud_blueprint_config.kms_custom_form.variables["namespace"] 
  }
}

data "autocloud_blueprint_config" "final" {
  source = {
    kms        = data.autocloud_blueprint_config.kms_custom_form.blueprint_config,
    s3         = data.autocloud_blueprint_config.s3_custom_form.blueprint_config
  }
}

resource "autocloud_blueprint" "example" {
  name = "KMS and S3 (reference from bc)"

  ###
  # UI Configuration
  #
  author       = "enrique.enciso@autocloud.dev"
  description  = "Terraform Generator KMS file vars"
  instructions = "Instructions..."

  labels = ["aws"]

  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new EKS generator (testing file vars), created by {{authorName}}"
      commit_message_template = "[AutoCloud] new EKS generator (testing file vars), created by {{authorName}}"
      body                    = file("../files/pull_request.md.tpl")
      variables = {
        authorName  = "generic.authorName"
        clusterName = "kms.name"
      }
    }
  }


  ###
  # File definitions
  #
  file {
    action      = "CREATE"
    destination = "aws/{{environment}}/{{namespace}}_{{environment}}_{{name}}_form-var-values.tf"

    variables = {
      // reference to kms_custom_form blueprint config    
      namespace   = data.autocloud_blueprint_config.kms_custom_form.variables["namespace"]
      environment = data.autocloud_blueprint_config.kms_custom_form.variables["environment"]
      name        = data.autocloud_blueprint_config.kms_custom_form.variables["name"]
      deletion     = data.autocloud_blueprint_config.kms_custom_form.variables["deletion_window_in_days"]
    }

    modules = [
      autocloud_module.kms.name
    ]
  }
  config = data.autocloud_blueprint_config.final.config
}


output "bf" {
  value = data.autocloud_blueprint_config.kms_custom_form.config
}
