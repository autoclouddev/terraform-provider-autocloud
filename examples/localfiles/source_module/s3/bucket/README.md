Terraform AWS S3 Bucket Module
==============================

## Overview

Creates a secure S3 bucket

## Specifications

#### Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement_terraform) | >= 0.14 |

#### Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider_aws) | n/a |

#### Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_default_label"></a> [default_label](#module_default_label) | ../../../../general/label | n/a |

#### Resources

| Name | Type |
|------|------|
| [aws_s3_bucket.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket) | resource |
| [aws_s3_bucket_public_access_block.default](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_public_access_block) | resource |
| [aws_iam_policy_document.default](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |

#### Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_account"></a> [account](#input_account) | Name of the account being used (master, nonprod, prod, etc) | `string` | `null` | no |
| <a name="input_attributes"></a> [attributes](#input_attributes) | Additional attributes (e.g. `1`) | `list(string)` | `[]` | no |
| <a name="input_authorized_users"></a> [authorized_users](#input_authorized_users) | A list of IAM user ARNs for users that are allowed to modify this S3 bucket | `list(string)` | `[]` | no |
| <a name="input_block_public"></a> [block_public](#input_block_public) | Toggle to block or allow public access | `bool` | `true` | no |
| <a name="input_cloud_provider"></a> [cloud_provider](#input_cloud_provider) | Cloud provider name if any | `string` | `null` | no |
| <a name="input_context"></a> [context](#input_context) | Single object for setting entire context at once.<br>See description of individual variables for details.<br>Leave string and numeric variables as `null` to use default value.<br>Individual variable settings (non-null) override settings in context object,<br>except for attributes, tags, and additional_tag_map, which are merged. | <pre>object({<br>    enabled             = bool<br>    namespace           = string<br>    cloud_provider      = string<br>    account             = string<br>    region              = string<br>    environment         = string<br>    stage               = string<br>    name                = string<br>    delimiter           = string<br>    attributes          = list(string)<br>    tags                = map(string)<br>    additional_tag_map  = map(string)<br>    regex_replace_chars = string<br>    label_order         = list(string)<br>    id_length_limit     = number<br>  })</pre> | <pre>{<br>  "account": null,<br>  "additional_tag_map": {},<br>  "attributes": [],<br>  "cloud_provider": null,<br>  "delimiter": null,<br>  "enabled": true,<br>  "environment": null,<br>  "id_length_limit": null,<br>  "label_order": [],<br>  "name": null,<br>  "namespace": null,<br>  "regex_replace_chars": null,<br>  "region": null,<br>  "stage": null,<br>  "tags": {}<br>}</pre> | no |
| <a name="input_domain"></a> [domain](#input_domain) | TLD to use when deploying assets | `string` | `null` | no |
| <a name="input_enabled"></a> [enabled](#input_enabled) | Set to false to prevent the module from creating any resources | `bool` | `true` | no |
| <a name="input_environment"></a> [environment](#input_environment) | Environment, e.g. 'prod', 'staging', 'dev', 'pre-prod', 'UAT' | `string` | `null` | no |
| <a name="input_error_file"></a> [error_file](#input_error_file) | Name of the error file to use if web app | `string` | `null` | no |
| <a name="input_external_policy"></a> [external_policy](#input_external_policy) | Whether or not to use an external S3 policy, configured and deployed after the bucket is created | `bool` | `false` | no |
| <a name="input_force_destroy"></a> [force_destroy](#input_force_destroy) | Force deletion of all contents on the deletion of bucket | `bool` | `false` | no |
| <a name="input_index_file"></a> [index_file](#input_index_file) | Name of the index file to use if web app | `string` | `null` | no |
| <a name="input_is_web_app"></a> [is_web_app](#input_is_web_app) | Is the bucket being used for a web application? | `bool` | `false` | no |
| <a name="input_kms_key_arn"></a> [kms_key_arn](#input_kms_key_arn) | KMS Key ARN to use encrypted the AMI. If omitted, defaults to AES256 | `string` | `""` | no |
| <a name="input_name"></a> [name](#input_name) | Module name, e.g. 'app' or 'jenkins' | `string` | `null` | no |
| <a name="input_namespace"></a> [namespace](#input_namespace) | Namespace, which could be your organization name or abbreviation, e.g. 'eg' or 'cp' | `string` | `null` | no |
| <a name="input_policy_statements"></a> [policy_statements](#input_policy_statements) | Additional policy statements for bucket policy. Must be defined in the module. See local variables for supported options | `list(string)` | <pre>[<br>  "DenyIncorrectEncryptionHeader",<br>  "DenyUnEncryptedObjectUploads"<br>]</pre> | no |
| <a name="input_region"></a> [region](#input_region) | AWS region to deploy asset into | `string` | `null` | no |
| <a name="input_stage"></a> [stage](#input_stage) | Stage, e.g. 'prod', 'staging', 'dev', OR 'source', 'build', 'test', 'deploy', 'release' | `string` | `null` | no |
| <a name="input_tags"></a> [tags](#input_tags) | Additional tags (e.g. `map('BusinessUnit','XYZ')` | `map(string)` | `{}` | no |

#### Outputs

| Name | Description |
|------|-------------|
| <a name="output_arn"></a> [arn](#output_arn) | AWS ARN of the bucket |
| <a name="output_bucket_domain_name"></a> [bucket_domain_name](#output_bucket_domain_name) | The bucket domain name. Will be of format bucketname.s3.amazonaws.com. |
| <a name="output_bucket_regional_domain_name"></a> [bucket_regional_domain_name](#output_bucket_regional_domain_name) | The bucket region-specific domain name. The bucket domain name including the region name, please refer here for format. Note: The AWS CloudFront allows specifying S3 region-specific endpoint when creating S3 origin, it will prevent redirect issues from CloudFront to S3 Origin URL. |
| <a name="output_context"></a> [context](#output_context) | Default label context |
| <a name="output_enabled"></a> [enabled](#output_enabled) | Whether or not the module is enabled |
| <a name="output_id"></a> [id](#output_id) | The name of the bucket |
| <a name="output_name"></a> [name](#output_name) | Name of the asset |
| <a name="output_policy"></a> [policy](#output_policy) | Original bucket policy for the bucket, to be modified and extended by external resources like Cloudfront |
| <a name="output_region"></a> [region](#output_region) | The AWS region this bucket resides in |

<!-- BEGIN_TF_DOCS -->
Terraform AWS S3 Bucket Module
==============================

## Overview

Creates a secure S3 bucket

## Specifications

#### Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement_aws) | ~> 4.0 |

#### Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider_aws) | ~> 4.0 |

#### Modules

No modules.

#### Resources

| Name | Type |
|------|------|
| [aws_s3_bucket.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket) | resource |
| [aws_s3_bucket_acl.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_acl) | resource |
| [aws_s3_bucket_cors_configuration.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_cors_configuration) | resource |
| [aws_s3_bucket_lifecycle_configuration.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_lifecycle_configuration) | resource |
| [aws_s3_bucket_ownership_controls.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_ownership_controls) | resource |
| [aws_s3_bucket_policy.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_policy) | resource |
| [aws_s3_bucket_public_access_block.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_public_access_block) | resource |
| [aws_s3_bucket_server_side_encryption_configuration.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_server_side_encryption_configuration) | resource |
| [aws_s3_bucket_versioning.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_versioning) | resource |
| [aws_s3_bucket_website_configuration.example](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_website_configuration) | resource |
| [aws_iam_policy_document.default](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |

#### Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_policy_template_filename"></a> [additional_policy_template_filename](#input_additional_policy_template_filename) | The path of a file containing a JSON document to add to the bucket policy as a source json | `string` | `null` | no |
| <a name="input_additional_policy_vars"></a> [additional_policy_vars](#input_additional_policy_vars) | A map of template variable inputs to populate the additional policy template with | `map(string)` | `{}` | no |
| <a name="input_authorized_users"></a> [authorized_users](#input_authorized_users) | A list of IAM user ARNs for users that are allowed to modify this S3 bucket | `list(string)` | `[]` | no |
| <a name="input_block_public"></a> [block_public](#input_block_public) | Toggle to block or allow public access | `bool` | `true` | no |
| <a name="input_cors_rules"></a> [cors_rules](#input_cors_rules) | A list of IAM user ARNs for users that are allowed to modify this S3 bucket | <pre>list(object({<br>    allowed_headers = optional(list(string))<br>    allowed_methods = optional(list(string))<br>    allowed_origins = optional(list(string))<br>    expose_headers  = optional(list(string))<br>    max_age_seconds = optional(number)<br>  }))</pre> | `[]` | no |
| <a name="input_enable_bucket_ownership_controls"></a> [enable_bucket_ownership_controls](#input_enable_bucket_ownership_controls) | Whether to use the bucket-ownership-control options for the bucket, default this is off, giving ObjectWriter ownership | `bool` | `false` | no |
| <a name="input_enabled"></a> [enabled](#input_enabled) | Set to false to prevent the module from creating any resources | `bool` | `true` | no |
| <a name="input_environment"></a> [environment](#input_environment) | Environment, e.g. 'prod', 'staging', 'dev', 'pre-prod', 'UAT' | `string` | `null` | no |
| <a name="input_error_file"></a> [error_file](#input_error_file) | Name of the error file to use if web app | `string` | `null` | no |
| <a name="input_external_policy"></a> [external_policy](#input_external_policy) | Whether or not to use an external S3 policy, configured and deployed after the bucket is created | `bool` | `false` | no |
| <a name="input_force_destroy"></a> [force_destroy](#input_force_destroy) | Force deletion of all contents on the deletion of bucket | `bool` | `false` | no |
| <a name="input_index_file"></a> [index_file](#input_index_file) | Name of the index file to use if web app | `string` | `null` | no |
| <a name="input_is_web_app"></a> [is_web_app](#input_is_web_app) | Is the bucket being used for a web application? | `bool` | `false` | no |
| <a name="input_kms_key_arn"></a> [kms_key_arn](#input_kms_key_arn) | KMS Key ARN to use encrypted the AMI. If omitted, defaults to AES256 | `string` | `""` | no |
| <a name="input_lifecycle_rules"></a> [lifecycle_rules](#input_lifecycle_rules) | A list of the lifecycle rules to apply to the S3 bucket | <pre>list(object({<br>    id      = string<br>    enabled = bool<br><br>    abort_incomplete_multipart_upload_days = optional(number)<br><br>    current_version_expiration = optional(object({<br>      date                         = optional(string)<br>      days                         = optional(number)<br>      expired_object_delete_marker = optional(bool)<br>    }))<br><br>    current_version_transition = optional(list(object({<br>      date          = optional(string)<br>      days          = optional(number)<br>      storage_class = string<br>    })))<br><br>    filter = optional(object({<br>      object_size_greater_than = optional(number)<br>      object_size_less_than    = optional(number)<br>      prefix                   = optional(string)<br>      tag = optional(object({<br>        key   = string<br>        value = string<br>      }))<br>    }))<br><br>    noncurrent_version_expiration = optional(object({<br>      newer_noncurrent_versions = optional(number)<br>      noncurrent_days           = optional(number)<br>    }))<br><br>    noncurrent_version_transition = optional(list(object({<br>      newer_noncurrent_versions = optional(number)<br>      noncurrent_days           = optional(number)<br>      storage_class             = string<br>    })))<br>  }))</pre> | `[]` | no |
| <a name="input_name"></a> [name](#input_name) | Module name, e.g. 'app' or 'jenkins' | `string` | `null` | no |
| <a name="input_namespace"></a> [namespace](#input_namespace) | Namespace, which could be your organization name or abbreviation, e.g. 'eg' or 'cp' | `string` | `null` | no |
| <a name="input_object_ownership"></a> [object_ownership](#input_object_ownership) | The object ownership rule applied to the bucket, default is ObjectWriter | `string` | `"ObjectWriter"` | no |
| <a name="input_policy_statements"></a> [policy_statements](#input_policy_statements) | Additional policy statements for bucket policy. Must be defined in the module. See local variables for supported options | `list(string)` | <pre>[<br>  "DenyIncorrectEncryptionHeader",<br>  "DenyUnEncryptedObjectUploads"<br>]</pre> | no |
| <a name="input_tags"></a> [tags](#input_tags) | Additional tags (e.g. `map('BusinessUnit','XYZ')` | `map(string)` | `{}` | no |

#### Outputs

| Name | Description |
|------|-------------|
| <a name="output_arn"></a> [arn](#output_arn) | AWS ARN of the bucket |
| <a name="output_bucket_domain_name"></a> [bucket_domain_name](#output_bucket_domain_name) | The bucket domain name. Will be of format bucketname.s3.amazonaws.com. |
| <a name="output_bucket_name"></a> [bucket_name](#output_bucket_name) | The bucket name. |
| <a name="output_bucket_regional_domain_name"></a> [bucket_regional_domain_name](#output_bucket_regional_domain_name) | The bucket region-specific domain name. The bucket domain name including the region name, please refer here for format. Note: The AWS CloudFront allows specifying S3 region-specific endpoint when creating S3 origin, it will prevent redirect issues from CloudFront to S3 Origin URL. |
| <a name="output_enabled"></a> [enabled](#output_enabled) | Whether or not the module is enabled |
| <a name="output_environment"></a> [environment](#output_environment) | Environment of the asset |
| <a name="output_error_file"></a> [error_file](#output_error_file) | Name of the error file to use if web app |
| <a name="output_id"></a> [id](#output_id) | The name of the bucket |
| <a name="output_index_file"></a> [index_file](#output_index_file) | Name of the index file to use if web app |
| <a name="output_name"></a> [name](#output_name) | Name of the asset |
| <a name="output_namespace"></a> [namespace](#output_namespace) | Namespace of the asset |
| <a name="output_ownership_control"></a> [ownership_control](#output_ownership_control) | The ownership control setting of the bucket. Will default to ObjectWriter if not specified or the resource is not used |
| <a name="output_policy"></a> [policy](#output_policy) | Original bucket policy for the bucket, to be modified and extended by external resources like Cloudfront |
| <a name="output_region"></a> [region](#output_region) | The AWS region this bucket resides in |
| <a name="output_tags"></a> [tags](#output_tags) | Tags for the asset |
<!-- END_TF_DOCS -->
