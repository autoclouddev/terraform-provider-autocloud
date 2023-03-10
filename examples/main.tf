terraform {
  required_providers {
    autocloud = {
      version = "0.2.0"
      source  = "autocloud.io/autoclouddev/autocloud"
    }
  }
}


/*
USE THIS FILE AS YOU NEED FIT, THIS IS JUST A PLAYGROUND

*/
provider "autocloud" {
  # endpoint = "https://api.autocloud.domain.com/api/v.0.0.1"
}

data "autocloud_github_repos" "repos" {}

####
# Local vara
locals {
  # Destination repos where generated code will be submitted
  dest_repos = [
    for repo in data.autocloud_github_repos.repos.data[*].url : repo
    if length(regexall("/infrastructure-live", repo)) > 0 || length(regexall("/self-hosted-infrastructure-live", repo)) > 0
  ]
}


locals {
  s3_vars_from_form_config = jsondecode(data.autocloud_blueprint_config.s3_processor.form_config)
  s3_vars_extend           = jsondecode(templatefile("${path.module}/files/s3bucket.vars.tpl", {}))
  # iterate over lists
  s3_vars_from_form_config_dict = { for item in local.s3_vars_from_form_config : item.id => item }
  s3_vars_extend_dict           = { for item in local.s3_vars_extend : item.id => item }
  # combine both dictionaries
  s3_form_config = jsonencode(values(merge(local.s3_vars_extend_dict, local.s3_vars_from_form_config_dict)))
}




resource "autocloud_module" "s3_bucket" {

  ####
  # Name of the generator
  name = "S3Bucket"

  ####
  # Can be any supported terraform source reference, must optionaly take version
  #
  #   source = "app.terraform.io/autocloud/aws/s3_bucket"
  #   version = "0.24.0"
  #
  # See docs: https://developer.hashicorp.com/terraform/language/modules/sources

  version = "3.4.0"
  source  = "terraform-aws-modules/s3-bucket/aws" // https://registry.terraform.io/modules/terraform-aws-modules/s3-bucket/aws/latest
  #tags_variable = "custom_tags"                         # this should not be a valid property
  display_order = ["bucket"]
}


resource "autocloud_module" "cloudfront" {

  ####
  # Name of the generator
  name = "Cloudfront"

  ####
  # Can be any supported terraform source reference, must optionaly take version
  #
  #   source = "app.terraform.io/autocloud/aws/s3_bucket"
  #   version = "0.24.0"
  #
  # See docs: https://developer.hashicorp.com/terraform/language/modules/sources

  version       = "3.0.0"
  source        = "terraform-aws-modules/cloudfront/aws" // https://registry.terraform.io/modules/terraform-aws-modules/cloudfront/aws/latest
  display_order = ["web_acl_id", "price_class", "http_version"]
}

data "autocloud_module" "s3_bucket_alt" {

  name = "S3BucketAlt"
  filter {
    name    = "S3Bucket"
    version = "3.4.0"
    source  = "terraform-aws-modules/s3-bucket/aws" // https://registry.terraform.io/modules/terraform-aws-modules/s3-bucket/aws/latest
  }

}

data "autocloud_blueprint_config" "s3_processor" {
  source_module_id = autocloud_module.s3_bucket.id
  omit_variables = [
    "acl",
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
    "object_ownership",
    "policy",
    "putin_khuylo",
    "request_payer",
  ]

  # bucket_prefix, acceleration_status => these vars are of 'shortText' type
  # attach_public_policy is of 'radio' type ('checkbox' types are similar to 'radio' types)

  # OVERRIDE VARIABLE EXAMPLES
  # - overriding bucket_prefix 'shortText' into 'radio'
  variable {
    name         = "bucket_prefix"
    display_name = "bucket prefix (from override block)"
    helper_text  = "bucket prefix helper text (from override block)"
    form_config {
      type = "radio"
      options {
        option {
          label   = "dev"
          value   = "some-dev-prefix"
          checked = false
        }
        option {
          label   = "nonprod"
          value   = "some-nonprod-prefix"
          checked = true
        }
        option {
          label = "prod"
          value = "some-prod-prefix"
        }
      }
      validation_rule {
        rule          = "isRequired"
        error_message = "invalid"
      }
    }
  }

  # - overriding acceleration_status 'shortText' into 'checkbox'
  variable {
    name = "acceleration_status"
    form_config {
      type = "checkbox"
      options {
        option {
          label = "Option 1"
          value = "acceleration_status_1"

        }
        option {
          label   = "Option 2"
          value   = "acceleration_status_2"
          checked = true
        }
        option {
          label   = "Option 3"
          value   = "acceleration_status_3"
          checked = true
        }
      }
    }
  }



  # - overriding bucket to make it mandatory
  variable {
    name         = "bucket"
    display_name = "Bucket name"
    helper_text  = "Set the bucket name"

    form_config {
      type = "shortText"

      validation_rule {
        rule          = "isRequired"
        error_message = "invalid"
      }
    }
  }
}
data "autocloud_blueprint_config" "cf_processor" {
  source_module_id = autocloud_module.cloudfront.id

  # omitting most of the variables to simplify the form
  omit_variables = [
    "aliases",
    // "comment", // we'll take the value from s3 bucket name
    "create_distribution",
    "create_monitoring_subscription",
    "create_origin_access_identity",
    "custom_error_response",
    "default_cache_behavior",
    // "default_root_object", // this will contain a default value
    // "enabled", // we'll let the user select this value in the form
    "geo_restriction",
    "http_version",
    "is_ipv6_enabled",
    "logging_config",
    "ordered_cache_behavior",
    "origin",
    "origin_access_identities",
    "origin_group",
    "price_class",
    "realtime_metrics_subscription_status",
    "retain_on_delete",
    "tags",
    "viewer_certificate",
    "wait_for_deployment",
    "web_acl_id",
  ]

  # OVERRIDE VARIABLE EXAMPLES
  # - set values from other modules outputs
  variable {
    name  = "comment"
    value = autocloud_module.s3_bucket.outputs["s3_bucket_id"]
  }

  variable {
    name  = "default_root_object"
    value = "index.html"
  }
}

resource "autocloud_blueprint" "example" {
  name = "S3andCloudfront"

  ###
  # UI Configuration
  #
  author       = "enrique.enciso@autocloud.dev"
  description  = "Terraform Generator storage in cloud"
  instructions = <<-EOT
  To deploy this generator, follow these simple steps:

  step 1: step-1-description
  step 2: step-2-description
  step 3: step-3-description
  EOT

  labels = ["aws"]

  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new EKS generator, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new EKS generator, created by {{authorName}}"
      body                    = file("./files/pull_request.md.tpl")
      variables = {
        authorName  = "generic.authorName"
        clusterName = "S3Bucket.Bucket"
        dummyParam  = autocloud_module.s3_bucket.variables["restrict_public_buckets"]
      }
    }
  }


  ###
  # File definitions
  #
  file {
    action = "CREATE"

    destination = "{{env}}/modules/s3_bucket_{{name}}.tf"

    # var name => variable id
    variables = {
      env  = "${autocloud_module.s3_bucket.name}.env"
      name = autocloud_module.s3_bucket.variables["bucket"]
    }

    modules = [autocloud_module.s3_bucket.name, autocloud_module.cloudfront.name]
  }

  ## these will change in future sprints, autocloud_module will be replaced
  autocloud_module {
    id = autocloud_module.s3_bucket.id

    # form_config = local.s3_form_config # example to append questions using locals
    form_config = data.autocloud_blueprint_config.s3_processor.form_config # config from data source
    # form_config = templatefile("${path.module}/files/s3bucket.vars.tpl", {})  # example from file
    # form_config = data.autocloud_form_config.s3_bucket.form_config            # example from data
    # form_config = autocloud_module.s3_bucket.blueprint_config                 # example from resource
    # form_config = data.autocloud_module.s3_bucket_alt.blueprint_config        # example from data

    #template_config = file("${path.module}/files/s3bucket.tf.tpl")
  }

  autocloud_module {
    id = autocloud_module.cloudfront.id

    form_config = data.autocloud_blueprint_config.cf_processor.form_config
  }
}
