# this is a template OVERRIDE
module "S3Bucket" {
    source                                = "terraform-aws-modules/s3-bucket/aws"
    version                               = "3.4.0"
    bucket                                = "{{bucket}}"
    bucket_prefix                         = "{{bucket_prefix}}"
    acl                                   = "{{acl}}"
    attach_deny_insecure_transport_policy = "{{attach_deny_insecure_transport_policy}}"
    attach_lb_log_delivery_policy         = "{{attach_lb_log_delivery_policy}}"
    force_destroy                         = "{{force_destroy}}"
    request_payer                         = "{{request_payer}}"
    object_ownership                      = "{{object_ownership}}"
    control_object_ownership              = "{{control_object_ownership}}"
    attach_public_policy                  = "{{attach_public_policy}}"
    block_public_policy                   = "{{block_public_policy}}"
    attach_policy                         = "{{attach_policy}}"
    acceleration_status                   = "{{acceleration_status}}"
    restrict_public_buckets               = "{{restrict_public_buckets}}"
    putin_khuylo                          = "{{putin_khuylo}}"
    expected_bucket_owner                 = "{{expected_bucket_owner}}"
    create_bucket                         = "{{create_bucket}}"
    block_public_acls                     = "{{block_public_acls}}"
    attach_elb_log_delivery_policy        = "{{attach_elb_log_delivery_policy}}"
    object_lock_enabled                   = "{{object_lock_enabled}}"
    attach_require_latest_tls_policy      = "{{attach_require_latest_tls_policy}}"
    ignore_public_acls                    = "{{ignore_public_acls}}"
    policy                                = "{{policy}}"
}