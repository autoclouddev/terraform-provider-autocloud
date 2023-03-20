# Local variables
locals {
  kms_key_arn = var.kms_key_arn == "" ? null : var.kms_key_arn

  default_policy_statements = [
    "DenyUnsecuredTransport"
  ]

  policy_statements = concat(local.default_policy_statements, var.policy_statements)

  create_additional_policy = var.additional_policy_template_filename != null

  bucket_name = join("-", [var.namespace, var.environment, var.name])

  tags = merge(
    {
      Environment = var.environment
      Name        = var.name
      Namespace   = var.namespace
    },
    var.tags
  )
}


####
# S3 bucket policy
#
# Define IAM policy for bucket
#
data "aws_iam_policy_document" "default" {
  count = var.enabled == true ? 1 : 0

  source_policy_documents = local.create_additional_policy ? [templatefile(var.additional_policy_template_filename, var.additional_policy_vars)] : null

  dynamic "statement" {
    for_each = contains(local.policy_statements, "DenyUnsecuredTransport") ? ["create"] : []

    content {
      sid = "DenyUnsecuredTransport"

      effect = "Deny"

      principals {
        type        = "*"
        identifiers = ["*"]
      }

      actions = [
        "s3:*"
      ]

      resources = [
        "arn:aws:s3:::${local.bucket_name}",
        "arn:aws:s3:::${local.bucket_name}/*",
      ]

      condition {
        test     = "Bool"
        variable = "aws:SecureTransport"

        values = [
          "false"
        ]
      }
    }
  }

  dynamic "statement" {
    for_each = contains(local.policy_statements, "DenyIncorrectEncryptionHeader") ? ["create"] : []

    content {
      sid = "DenyIncorrectEncryptionHeader"

      effect = "Deny"
      principals {
        type        = "*"
        identifiers = ["*"]
      }
      actions = [
        "s3:PutObject"
      ]

      resources = [
        "arn:aws:s3:::${local.bucket_name}",
        "arn:aws:s3:::${local.bucket_name}/*",
      ]

      condition {
        test     = "StringNotEquals"
        variable = "s3:x-amz-server-side-encryption"

        values = local.kms_key_arn != null ? ["aws:kms"] : ["AES256"]
      }
    }
  }


  dynamic "statement" {
    for_each = contains(local.policy_statements, "DenyUnEncryptedObjectUploads") ? ["create"] : []

    content {
      sid = "DenyUnEncryptedObjectUploads"

      effect = "Deny"
      principals {
        type        = "*"
        identifiers = ["*"]
      }
      actions = [
        "s3:PutObject",
      ]

      resources = [
        "arn:aws:s3:::${local.bucket_name}",
        "arn:aws:s3:::${local.bucket_name}/*",
      ]

      condition {
        test     = "Null"
        variable = "s3:x-amz-server-side-encryption"

        values = [
          "true"
        ]
      }
    }
  }

  dynamic "statement" {
    for_each = contains(local.policy_statements, "AllowAuthorizedUsers") && length(var.authorized_users) > 0 ? ["create"] : []

    content {
      sid = "AllowAuthorizedUsers"

      effect = "Allow"
      principals {
        type        = "AWS"
        identifiers = var.authorized_users
      }
      actions = [
        "s3:*"
      ]

      resources = [
        "arn:aws:s3:::${local.bucket_name}",
        "arn:aws:s3:::${local.bucket_name}/*",
      ]

      condition {
        test     = "Null"
        variable = "s3:x-amz-server-side-encryption"

        values = [
          "true"
        ]
      }
    }
  }
}



####
# S3 bucket
#
# Main bucket resource
#
resource "aws_s3_bucket" "this" {
  count  = var.enabled == true ? 1 : 0
  bucket = local.bucket_name

  force_destroy = var.force_destroy

  tags = local.tags

  # Terraform currently restricts ignore_changes to static values due to it's use in building dependency graphs.
  #
  # See here: https://github.com/hashicorp/terraform/issues/24188
  #
  # Proposed changes are being discussed to allow statements like the following:
  #
  # lifecycle {
  #   ignore_changes = var.external_policy ? ["policy"] : []
  # }
  #
  # However, this will be a long time coming. In the mean time, the Bucket module will rewrite any changes made
  # Outside of the stack. External policies will have to be rerun/recreated after this module executes.

  lifecycle {
    ignore_changes = [
      lifecycle_rule
    ]
  }
}



####
# Bucket Policy
#
# If the bucket policy is not marked as external, set the bucket policy
#
resource "aws_s3_bucket_policy" "this" {
  count = var.enabled == true && !var.external_policy ? 1 : 0

  bucket = join("", aws_s3_bucket.this[*].id)

  policy = join("", data.aws_iam_policy_document.default[*].json)
}



####
# Enable versioning
#
resource "aws_s3_bucket_versioning" "this" {
  count = var.enabled == true ? 1 : 0

  bucket = join("", aws_s3_bucket.this[*].id)

  versioning_configuration {
    status = "Enabled"
  }

  # Terraform reccomends waiting 15 min after enabling versioning to use the bucket due to internal
  # propagation delays.
  #
  # See https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_versioning
  provisioner "local-exec" {
    command = "echo 'Sleeping for 15 minutes to wait for versioning configuration propagation'; sleep 900"
  }
}



####
# Public access block
#
# Block all public access
#
resource "aws_s3_bucket_public_access_block" "this" {
  count = var.enabled == true && var.block_public == true ? 1 : 0

  bucket = join("", aws_s3_bucket.this[*].id)

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}



####
# Bucket ACL
#
# Make bucket private
# If BucketOwnership is defined as BucketOwnerEnforced, then ACLs are not used, only bucket policies!
resource "aws_s3_bucket_acl" "this" {
  count = var.enabled == true && var.object_ownership != "BucketOwnerEnforced" ? 1 : 0

  bucket = join("", aws_s3_bucket.this[*].id)
  acl    = "private"
}



####
# Enforce encryption
#
# Set SSE to either AES256 or KMS as appropriate depending on whether or not a KMS key ARN was provided
#
resource "aws_s3_bucket_server_side_encryption_configuration" "this" {
  count = var.enabled == true && var.block_public == true ? 1 : 0

  bucket = join("", aws_s3_bucket.this[*].id)

  rule {
    dynamic "apply_server_side_encryption_by_default" {
      for_each = local.kms_key_arn != null ? [local.kms_key_arn] : []
      iterator = kms_key_arn

      content {
        kms_master_key_id = kms_key_arn.value
        sse_algorithm     = "aws:kms"
      }
    }

    dynamic "apply_server_side_encryption_by_default" {
      for_each = local.kms_key_arn == null ? ["AES256"] : []
      iterator = algorithm

      content {
        sse_algorithm = algorithm.value
      }
    }
  }
}



####
# CORS configuration
#
# Set CORS rules
#
resource "aws_s3_bucket_cors_configuration" "this" {
  count = var.enabled == true && length(var.cors_rules) > 0 ? 1 : 0

  bucket = join("", aws_s3_bucket.this[*].id)

  dynamic "cors_rule" {
    for_each = toset(var.cors_rules)
    iterator = cors

    content {
      allowed_headers = lookup(cors.value, "allowed_headers", null) # ["*"]
      allowed_methods = lookup(cors.value, "allowed_methods", null) # ["PUT", "POST"]
      allowed_origins = lookup(cors.value, "allowed_origins", null) # ["https://s3-website-test.hashicorp.com"]
      expose_headers  = lookup(cors.value, "expose_headers", null)  # ["ETag"]
      max_age_seconds = lookup(cors.value, "max_age_seconds", null) # 3000
    }
  }
}



####
# Website configuration
#
# Set bucket website configuration
#
resource "aws_s3_bucket_website_configuration" "example" {
  count = var.enabled && var.is_web_app == true ? 1 : 0

  bucket = join("", aws_s3_bucket.this[*].id)

  index_document {
    suffix = var.index_file != null ? var.index_file : "index.html"
  }

  error_document {
    key = var.error_file != null ? var.error_file : "error.html"
  }
}




resource "aws_s3_bucket_lifecycle_configuration" "this" {
  count = var.enabled && length(var.lifecycle_rules) > 0 ? 1 : 0

  bucket = join("", aws_s3_bucket.this[*].id)

  dynamic "rule" {
    for_each = toset(var.lifecycle_rules)

    content {
      id = rule.value.id

      status = rule.value.enabled ? "Enabled" : "Disabled"

      dynamic "abort_incomplete_multipart_upload" {
        for_each = try(rule.value.abort_incomplete_multipart_upload_days > 0, false) ? ["create"] : []
        content {
          days_after_initiation = rule.value.abort_incomplete_multipart_upload_days
        }
      }

      dynamic "expiration" {
        for_each = rule.value.current_version_expiration != null ? ["create"] : []

        content {
          date = lookup(rule.value.current_version_expiration, "date", null)
          days = lookup(rule.value.current_version_expiration, "days", null)
        }
      }

      dynamic "filter" {
        for_each = rule.value.filter != null ? ["create"] : []

        content {
          object_size_greater_than = lookup(rule.value.filter, "object_size_greater_than", null)
          object_size_less_than    = lookup(rule.value.filter, "object_size_less_than", null)
          prefix                   = lookup(rule.value.filter, "prefix", null)

          dynamic "tag" {
            for_each = rule.value.filter.tag != null ? ["create"] : []

            content {
              key   = lookup(rule.value.filter.tag, "key", null)
              value = lookup(rule.value.filter.tag, "value", null)
            }
          }
        }
      }

      dynamic "noncurrent_version_expiration" {
        for_each = rule.value.noncurrent_version_expiration != null ? ["create"] : []

        content {
          newer_noncurrent_versions = try(lookup(rule.value.noncurrent_version_expiration, "newer_noncurrent_versions", null), null)
          noncurrent_days           = try(lookup(rule.value.noncurrent_version_expiration, "noncurrent_days", null), null)
        }
      }

      dynamic "noncurrent_version_transition" {
        for_each = toset(
          can(rule.value["noncurrent_version_transition"]) && rule.value["noncurrent_version_transition"] != null ? rule.value["noncurrent_version_transition"] : []
        )

        content {
          newer_noncurrent_versions = can(noncurrent_version_transition.value["newer_noncurrent_versions"]) && noncurrent_version_transition.value["newer_noncurrent_versions"] != null ? noncurrent_version_transition.value["newer_noncurrent_versions"] : null
          noncurrent_days           = can(noncurrent_version_transition.value["noncurrent_days"]) && noncurrent_version_transition.value["noncurrent_days"] != null ? noncurrent_version_transition.value["noncurrent_days"] : null
          storage_class             = noncurrent_version_transition.value.storage_class
        }
      }

      dynamic "transition" {
        for_each = toset(
          can(rule.value["current_version_transition"]) && rule.value["current_version_transition"] != null ? rule.value["current_version_transition"] : []
        )

        content {
          date          = can(transition.value["date"]) && transition.value["date"] != null ? transition.value["date"] : null
          days          = can(transition.value["days"]) && transition.value["days"] != null ? transition.value["days"] : null
          storage_class = transition.value.storage_class
        }
      }
    }
  }
}


resource "aws_s3_bucket_ownership_controls" "this" {
  count  = var.enabled && var.enable_bucket_ownership_controls ? 1 : 0
  bucket = join("", aws_s3_bucket.this[*].id)
  rule {
    object_ownership = var.object_ownership
  }
}
