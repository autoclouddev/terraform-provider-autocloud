###
# Standard Variables
#
# All modules expect these variables to be defined. Omitting them will result in errors or undefined/undesirable behavior.
#
variable "account" {
  type        = string
  default     = null
  description = "Name of the account being used (master, nonprod, prod, etc)"
}

variable "attributes" {
  type        = list(string)
  default     = []
  description = "Additional attributes (e.g. `1`)"
}

variable "cloud_provider" {
  type        = string
  default     = null
  description = "Cloud provider name if any"
}

variable "context" {
  type = object({
    enabled             = bool
    namespace           = string
    cloud_provider      = string
    account             = string
    region              = string
    environment         = string
    stage               = string
    domain              = string
    name                = string
    delimiter           = string
    attributes          = list(string)
    tags                = map(string)
    additional_tag_map  = map(string)
    regex_replace_chars = string
    label_order         = list(string)
    id_length_limit     = number
  })
  default = {
    enabled             = true
    namespace           = null
    cloud_provider      = null
    account             = null
    region              = null
    environment         = null
    stage               = null
    domain              = null
    name                = null
    delimiter           = null
    attributes          = []
    tags                = {}
    additional_tag_map  = {}
    regex_replace_chars = null
    label_order         = []
    id_length_limit     = null
  }
  description = <<-EOT
    Single object for setting entire context at once.
    See description of individual variables for details.
    Leave string and numeric variables as `null` to use default value.
    Individual variable settings (non-null) override settings in context object,
    except for attributes, tags, and additional_tag_map, which are merged.
  EOT
}

variable "domain" {
  type        = string
  default     = null
  description = "TLD to use when deploying assets"
}

variable "enabled" {
  type        = bool
  default     = true
  description = "Set to false to prevent the module from creating any resources"
}

variable "environment" {
  type        = string
  default     = null
  description = "Environment, e.g. 'prod', 'staging', 'dev', 'pre-prod', 'UAT'"
}

variable "name" {
  type        = string
  default     = null
  description = "Module name, e.g. 'app' or 'jenkins'"
}

variable "namespace" {
  type        = string
  default     = null
  description = "Namespace, which could be your organization name or abbreviation, e.g. 'eg' or 'cp'"
}

variable "region" {
  type        = string
  default     = null
  description = "AWS region to deploy asset into"
}

variable "stage" {
  type        = string
  default     = null
  description = "Stage, e.g. 'prod', 'staging', 'dev', OR 'source', 'build', 'test', 'deploy', 'release'"
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Additional tags (e.g. `map('BusinessUnit','XYZ')`"
}

###
# Module Variables
#
# Variables specific to this module
#
variable "authorized_users" {
  type        = list(string)
  default     = []
  description = "A list of IAM user ARNs for users that are allowed to modify this S3 bucket"
}


variable "cors_rules" {
  type = list(object({
    allowed_headers = optional(list(string))
    allowed_methods = optional(list(string))
    allowed_origins = optional(list(string))
    expose_headers  = optional(list(string))
    max_age_seconds = optional(number)
  }))

  default     = []
  description = "A list of IAM user ARNs for users that are allowed to modify this S3 bucket"
}


variable "external_policy" {
  type        = bool
  default     = false
  description = "Whether or not to use an external S3 policy, configured and deployed after the bucket is created"
}

variable "kms_key_arn" {
  type        = string
  default     = ""
  description = "KMS Key ARN to use encrypted the AMI. If omitted, defaults to AES256"
}

variable "force_destroy" {
  type        = bool
  default     = false
  description = "Force deletion of all contents on the deletion of bucket"
}

variable "policy_statements" {
  type = list(string)
  default = [
    "DenyIncorrectEncryptionHeader",
    "DenyUnEncryptedObjectUploads"
  ]
  description = "Additional policy statements for bucket policy. Must be defined in the module. See local variables for supported options"
}

variable "is_web_app" {
  description = "Is the bucket being used for a web application?"
  default     = false
  type        = bool
}

variable "index_file" {
  description = "Name of the index file to use if web app"
  default     = null
  type        = string
}

variable "error_file" {
  description = "Name of the error file to use if web app"
  default     = null
  type        = string
}

variable "block_public" {
  description = "Toggle to block or allow public access"
  default     = true
  type        = bool
}

variable "additional_policy_template" {
  description = "A JSON document to add to the bucket policy as a source json"
  default     = null
  type        = string
}

variable "additional_policy_vars" {
  description = "A map of template variable inputs to populate the additional policy template with"
  default     = {}
  type        = map(string)
}

variable "lifecycle_rules" {
  description = "A list of the lifecycle rules to apply to the S3 bucket"
  default     = []
  type = list(object({
    id      = string
    enabled = bool

    abort_incomplete_multipart_upload_days = optional(number)

    expiration = optional(object({
      date                         = optional(string)
      days                         = optional(number)
      expired_object_delete_marker = optional(bool)
    }))

    filter = optional(object({
      object_size_greater_than = optional(number)
      object_size_less_than    = optional(number)
      prefix                   = optional(string)
      tag = optional(object({
        key   = string
        value = string
      }))
    }))

    noncurrent_version_expiration = optional(object({
      newer_noncurrent_versions = optional(number)
      noncurrent_days           = optional(number)
    }))

    noncurrent_version_transition = optional(object({
      newer_noncurrent_versions = optional(number)
      noncurrent_days           = optional(number)
      storage_class             = string
    }))

    transition = optional(list(object({
      date          = optional(string)
      days          = optional(number)
      storage_class = string
    })))
  }))
}

variable "enable_bucket_ownership_controls" {
  type        = bool
  default     = false
  description = "Whether to use the bucket-ownership-control options for the bucket, default this is off, giving ObjectWriter ownership"
}

variable "object_ownership" {
  type        = string
  default     = "ObjectWriter"
  description = "The object ownership rule applied to the bucket, default is ObjectWriter"
  validation {
    condition     = contains(["BucketOwnerPreferred", "ObjectWriter", "BucketOwnerEnforced"], var.object_ownership)
    error_message = "Only valid values of: BucketOwnerPreferred, ObjectWriter, BucketOwnerEnforced are supported."
  }
}
