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
    // "ecs_cluster",                       // object({ arn = string name = string })
    //"ecs_subnet_ids",                  // list(string)
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
    //"hello_world_container_ports",     // list(number)
    "lb_target_groups",                  // list( object({ container_port = number container_health_check_port = number lb_target_group_arn = string } ) )
    "logs_cloudwatch_group",             // string
    "logs_cloudwatch_retention",         // number
    "manage_ecs_security_group",         // bool
    "nlb_subnet_cidr_blocks",            // list(string)
    "service_registries",                // list(object({ registry_arn = string container_name = string container_port = number port = number }))
    "target_container_name",             // string
    "tasks_desired_count",               // number
    "tasks_maximum_percent",             // number
    "tasks_minimum_healthy_percent",     // number
  ]

  variable {
    name         = "ecs_cluster"
    required_values = jsonencode({
      "managed-by" = "autocloud"
      owner = null
      "business-unit" = [
        "finance",
        "legal",
        "engineering",
        "sales"
      ]
    })
  }

  variable {
    name         = "ecs_subnet_ids"
    display_name = "ecs subnet ids"
    helper_text  = "select the ecs subnet ids"
    type         = "checkbox"
    options {
      option {
        label   = "SUBNET 1"
        value   = "subnet-id-1"
        checked = true
      }
      option {
        label   = "SUBNET 2"
        value   = "subnet-id-2"
        checked = true
      }
      option {
        label   = "SUBNET 3"
        value   = "subnet-id-3"
        checked = false
      }
    }
  }

  variable {
    name         = "hello_world_container_ports"
    display_name = "hello work ports"
    helper_text  = "select the hello work ports"
    options {
      option {
        label   = "PORT 8080"
        value   = 8080
        checked = true
      }
      option {
        label   = "PORT 8081"
        value   = 8081
        checked = true
      }
      option {
        label   = "PORT 8082"
        value   = 8082
        checked = false
      }
    }
    required_values = jsonencode([
      "3002",
      "3001",
    ])
  }
}

resource "autocloud_blueprint" "example" {
  name = "ECS string and number lists (basic)"

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
      title                   = "[AutoCloud] new ECS generator, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new ECS generator, created by {{authorName}}"
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
