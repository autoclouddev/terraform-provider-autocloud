

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

resource "autocloud_module" "kms" {
  name   = "kms"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/security/kms?ref=0.1.0"
}

data "autocloud_blueprint_config" "kms_custom_form" {
  source = {
    kms = autocloud_module.kms.blueprint_config
  }

  omit_variables = [
    # kms variables
    "deletion_window_in_days",
    //"description",
    "enable_key_rotation",
    "enabled",
    //"environment",
    //"name",
    //"namespace",
  ]

  variable {
    name         = "environment"
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

  variable {
    name  = "name"
    value = "this-is-the-name"
  }

  variable {
    name  = "namespace"
    value = "this-is-the-namespace"
  }
}

resource "autocloud_blueprint" "example" {
  name = "KMS (file vars)"

  ###
  # UI Configuration
  #
  author       = "enrique.enciso@autocloud.dev"
  description  = "Terraform Generator KMS file vars"
  instructions = "Instructions..."

  labels = ["aws"]

  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new EKS generator (testing file vars), created by {{authorName}}"
      commit_message_template = "[AutoCloud] new EKS generator (testing file vars), created by {{authorName}}"
      body                    = file("../files/pull_request.md.tpl")
      variables = {
        authorName  = "generic.authorName"
        clusterName = "kms.name"
      }
    }
  }


  ###
  # File definitions
  #

  file {
    action      = "CREATE"
    destination = "aws/{{environment}}/{{namespace}}_{{environment}}_{{name}}_hardcoded-var-values.tf"
    variables = {
      namespace   = "unstyl"
      environment = "nonprod"
      name        = "static-generator"
    }

    modules = [
      autocloud_module.kms.name
    ]
  }
  file {
    action      = "CREATE"
    destination = "aws/{{environment}}/{{namespace}}_{{environment}}_{{name}}_form-var-values.tf"
    variables = {
      namespace   = "kms.namespace" // referencing the variable value using its name (note: it will break if the autocloud_module's name [kms] changes and this value is not updated)
      environment = autocloud_module.kms.variables.environment
      name        = autocloud_module.kms.variables["name"] // getting the <module>.<variable-name> format from its variable computed field
    }

    modules = [
      autocloud_module.kms.name
    ]
  }
  config = data.autocloud_blueprint_config.kms_custom_form.config
}


output "bf" {
  value = data.autocloud_blueprint_config.kms_custom_form.config
}
