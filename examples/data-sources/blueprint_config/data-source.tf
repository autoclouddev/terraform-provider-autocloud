terraform {
  required_providers {
    autocloud = {
      source  = "autoclouddev/autocloud"
      version = "~>0.9"
    }
  }
}

resource "autocloud_module" "s3_bucket" {
  name   = "s3bucket"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/storage/s3/bucket?ref=0.10.2"
}

data "autocloud_blueprint_config" "example" {
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
    "tags",
    "s3.variables.enabled"
  ]

  variable {
    name         = "environment"
    display_name = "environment"
    helper_text  = "select the environment"
    type         = "raw"
    value        = "var.ami"
  }
}
