# To show how to set override default values

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

  sample_string_var = "locals-subnet-1"
}

resource "autocloud_module" "s3_bucket" {
  name    = "S3Bucket"
  version = "3.4.0"
  source  = "terraform-aws-modules/s3-bucket/aws"
}

// https://registry.terraform.io/modules/trussworks/ecs-service/aws/latest
resource "autocloud_module" "ecs" {
  name    = "ecs"
  source  = "trussworks/ecs-service/aws"
  version = "6.6.0"
}

data "autocloud_blueprint_config" "generic" {
  variable {
    name         = "env"
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
}

data "autocloud_blueprint_config" "ecs_custom_form" {
  source = {
    ecs     = autocloud_module.ecs.blueprint_config
    generic = data.autocloud_blueprint_config.generic.blueprint_config
  }

  omit_variables = [
    //"ecs_cluster", // object({ arn = string name = string })
    //"ecs_subnet_ids",                  // list(string)
    "ecs_vpc_id",  // string
    "environment", // string
    "kms_key_id",  // string
    //"name",                              // string
    "additional_security_group_ids", // list(string)
    "alb_security_group",            // string
    //"assign_public_ip",                  // bool
    "associate_alb",                  // bool
    "associate_nlb",                  // bool
    "cloudwatch_alarm_actions",       // list(string)
    "cloudwatch_alarm_cpu_enable",    // bool
    "cloudwatch_alarm_cpu_threshold", // number
    "cloudwatch_alarm_mem_enable",    // bool
    "cloudwatch_alarm_mem_threshold", // number
    "cloudwatch_alarm_name",          // string
    "container_definitions",          // string
    "container_image",                // string
    //"container_volumes",                 // list( object({ name = string }) )
    "ec2_create_task_execution_role",    // bool
    "ecr_repo_arns",                     // list(string)
    "ecs_exec_enable",                   // bool
    "ecs_instance_role",                 // string
    "ecs_use_fargate",                   // bool
    "fargate_platform_version",          // string
    "fargate_task_cpu",                  // number
    "fargate_task_memory",               // number
    "health_check_grace_period_seconds", // number
    //"hello_world_container_ports",     // list(number)
    //"lb_target_groups",          // list( object({ container_port = number container_health_check_port = number lb_target_group_arn = string } ) )
    "logs_cloudwatch_group",     // string
    "logs_cloudwatch_retention", // number
    "manage_ecs_security_group", // bool
    "nlb_subnet_cidr_blocks",    // list(string)
    //"service_registries",        // list(object({ registry_arn = string container_name = string container_port = number port = number }))
    "target_container_name", // string
    //"tasks_desired_count",           // number
    "tasks_maximum_percent",         // number
    "tasks_minimum_healthy_percent", // number
  ]

  # string
  variable {
    name  = "name"
    value = "dummy-name"
  }

  # bool
  variable {
    name  = "assign_public_ip"
    value = true
  }

  # number
  variable {
    name  = "tasks_desired_count"
    value = 1000
  }

  # list(string)
  variable {
    name  = "ecs_subnet_ids"
    value = jsonencode(["subnet1", "subnet2", autocloud_module.s3_bucket.outputs["s3_bucket_id"]])
  }

  # list(number)
  variable {
    name  = "hello_world_container_ports"
    value = jsonencode([10000, 20000, 30000])
  }

  # object - object({ arn = string name = string })
  variable {
    name            = "ecs_cluster"
    required_values = jsonencode({ arn = "dummy-ecs-cluster-arn-value", name = "dummy-ecs-cluster-arn-name" })
  }

  // TODO: Support for references values
  # list(object) - list( object({ name = string }) )
  variable {
    name = "container_volumes"
    required_values = jsonencode([{ name = "dummy-volume-1" }, { name = "dummy-volume-2" }, { name = autocloud_module.s3_bucket.outputs["s3_bucket_id"] }])
  }

  # list(object) - list( object({ container_port = number container_health_check_port = number lb_target_group_arn = string } ) )
  variable {
    name = "lb_target_groups"
    required_values = jsonencode([
      { container_port = "10000", container_health_check_port = "10001", lb_target_group_arn = "dummy-arn-1" },
      { container_port = "20000", container_health_check_port = "20001", lb_target_group_arn = autocloud_module.s3_bucket.outputs["s3_bucket_id"] }
    ])
  }

  # list(object) - list(object({ registry_arn = string container_name = string container_port = number port = number }))
  variable {
    name = "service_registries"
    value = jsonencode([
      { registry_arn = "dummy-registry-arn-1", container_name = "dummy-container-name-1", container_port = "2000", port = "2001" },
      { registry_arn = "dummy-registry-arn-2", container_name = "dummy-container-name-2", container_port = "2000", port = autocloud_module.s3_bucket.outputs["s3_bucket_id"] }
    ])
  }

}

resource "autocloud_blueprint" "example" {
  name = "ECS override with static values"

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
      title                   = "[AutoCloud] new ECS generator (override with static values), created by {{authorName}}"
      commit_message_template = "[AutoCloud] new ECS generator (override with static values), created by {{authorName}}"
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
    destination = "ecs.tf"
    variables = {
    }

    modules = [
      autocloud_module.ecs.name
    ]
  }
  config = data.autocloud_blueprint_config.ecs_custom_form.config
}
