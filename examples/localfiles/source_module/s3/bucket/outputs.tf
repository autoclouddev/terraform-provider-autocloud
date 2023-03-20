###
# Standard Outputs
#
output "enabled" {
  description = "Whether or not the module is enabled"
  value       = var.enabled
}

output "environment" {
  description = "Environment of the asset"
  value       = var.environment
}

output "name" {
  description = "Name of the asset"
  value       = var.name
}

output "namespace" {
  description = "Namespace of the asset"
  value       = var.namespace
}

output "tags" {
  description = "Tags for the asset"
  value       = var.tags
}

###
# Module Outputs
#
output "arn" {
  value       = join("", aws_s3_bucket.this[*].arn)
  description = "AWS ARN of the bucket"
}

output "bucket_name" {
  value       = join("", aws_s3_bucket.this[*].id)
  description = "The bucket name."
}

output "bucket_domain_name" {
  value       = join("", aws_s3_bucket.this[*].bucket_domain_name)
  description = "The bucket domain name. Will be of format bucketname.s3.amazonaws.com."
}

output "bucket_regional_domain_name" {
  value       = join("", aws_s3_bucket.this[*].bucket_regional_domain_name)
  description = "The bucket region-specific domain name. The bucket domain name including the region name, please refer here for format. Note: The AWS CloudFront allows specifying S3 region-specific endpoint when creating S3 origin, it will prevent redirect issues from CloudFront to S3 Origin URL."
}

output "error_file" {
  value       = var.error_file
  description = "Name of the error file to use if web app"
}

output "id" {
  value       = join("", aws_s3_bucket.this[*].id)
  description = "The name of the bucket"
}

output "index_file" {
  value       = var.index_file
  description = "Name of the index file to use if web app"
}

output "ownership_control" {
  value       = var.object_ownership
  description = "The ownership control setting of the bucket. Will default to ObjectWriter if not specified or the resource is not used"
}

output "policy" {
  value       = join("", data.aws_iam_policy_document.default[*].json)
  description = "Original bucket policy for the bucket, to be modified and extended by external resources like Cloudfront"
  sensitive   = true
}

output "region" {
  value       = join("", aws_s3_bucket.this[*].region)
  description = "The AWS region this bucket resides in"
}
