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

// https://registry.terraform.io/modules/trussworks/ecs-service/aws/latest
resource "autocloud_module" "ecs" {
  name    = "ecs"
  source  = "trussworks/ecs-service/aws"
  version = "6.6.0"
}

data "autocloud_blueprint_config" "ecs_custom_form" {
  source = {
    ecs = autocloud_module.ecs.blueprint_config
  }

  omit_variables = [
    "ecs_cluster",                       // object({ arn = string name = string })
    "ecs_subnet_ids",                  // list(string)
    "ecs_vpc_id",                        // string
    "environment",                       // string
    "kms_key_id",                        // string
    "name",                              // string
    "additional_security_group_ids",     // list(string)
    "alb_security_group",                // string
    "assign_public_ip",                  // bool
    "associate_alb",                     // bool
    "associate_nlb",                     // bool
    "cloudwatch_alarm_actions",          // list(string)
    "cloudwatch_alarm_cpu_enable",       // bool
    "cloudwatch_alarm_cpu_threshold",    // number
    "cloudwatch_alarm_mem_enable",       // bool
    "cloudwatch_alarm_mem_threshold",    // number
    "cloudwatch_alarm_name",             // string
    "container_definitions",             // string
    "container_image",                   // string
    "container_volumes",                 // list( object({ name = string }) )
    "ec2_create_task_execution_role",    // bool
    "ecr_repo_arns",                     // list(string)
    "ecs_exec_enable",                   // bool
    "ecs_instance_role",                 // string
    "ecs_use_fargate",                   // bool
    "fargate_platform_version",          // string
    "fargate_task_cpu",                  // number
    "fargate_task_memory",               // number
    "health_check_grace_period_seconds", // number
    "hello_world_container_ports",     // list(number)
    //"lb_target_groups",                  // list( object({ container_port = number container_health_check_port = number lb_target_group_arn = string } ) )
    "logs_cloudwatch_group",             // string
    "logs_cloudwatch_retention",         // number
    "manage_ecs_security_group",         // bool
    "nlb_subnet_cidr_blocks",            // list(string)
    //"service_registries",                // list(object({ registry_arn = string container_name = string container_port = number port = number }))
    "target_container_name",             // string
    "tasks_desired_count",               // number
    "tasks_maximum_percent",             // number
    "tasks_minimum_healthy_percent",     // number
  ]
}


# ####
# # CloudFront Distribution
# #
# # Serves the content to the public internet and configures access to the private S3 bucket
resource "autocloud_module" "cloudfront" {
  name   = "cloudfront"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/network/cloudfront/distro?ref=0.4.0"
}

data "autocloud_blueprint_config" "cloudfront_processor" {
  source = {
    cloudfront = autocloud_module.cloudfront.blueprint_config
    ecs = data.autocloud_blueprint_config.ecs_custom_form.blueprint_config
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

  # list(object) Static value
  variable {
    name = "lambda_functions"
    value = jsonencode([
      { event_type = "delete", include_body = false, lambda_arn = "abc:123:lambda" },
      { event_type = "create", include_body = false, lambda_arn = "def:456:lambda" }
    ])
  }
}

####
# Create Blueprint
#
# Create generator blueprint that contains all the elements
resource "autocloud_blueprint" "this" {
  name = "Generator for List Objects"

  ###
  # UI Configuration
  #
  author       = "marco.franceschi@autocloud.dev"
  description  = "List objects generator"
  instructions = <<-EOT
  To deploy this generator, follow these simple steps:

  step 1: step-1-description
  step 2: step-2-description
  step 3: step-3-description
  EOT

  labels = ["aws"]

  ###
  # Form configuration
  config = data.autocloud_blueprint_config.cloudfront_processor.config


  ###
  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new List Objects, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new List Objects, created by {{authorName}}"
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
    destination = "aws/{{environment}}/{{namespace}}-{{environment}}-{{name}}.tf"
    variables = {
      namespace   = "unstyl"
      environment = "nonprod"
      name        = "static-generator"
    }

    modules = [
      autocloud_module.cloudfront.name,
    ]
  }
}
