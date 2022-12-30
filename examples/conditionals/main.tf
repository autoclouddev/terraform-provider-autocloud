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
# Local vara
locals {
  # Destination repos where generated code will be submitted
  dest_repos = [
    for repo in data.autocloud_github_repos.repos.data[*].url : repo
    if length(regexall("/terraform-generator-test", repo)) > 0 || length(regexall("/infrastructure-live-demo", repo)) > 0 || length(regexall("/self-hosted-infrastructure-live", repo)) > 0
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
    "environment",
    //"name",
    //"namespace",
  ]

  variable {
    name         = "name"
    display_name = "choose the usage of this kms (name)"
    helper_text  = "select the corresponding team"
    form_config {
      type = "radio"
      options {
        option {
          label   = "engineering"
          value   = "engineering"
          checked = false
        }
        option {
          label   = "finances"
          value   = "finances"
          checked = false
        }
        option {
          label   = "sales"
          value   = "sales"
          checked = true
        }
      }
    }
  }

  # ### overriding "namespace" (string)
  variable {
    name = "namespace"

    conditional {
      source    = "kms.name"
      condition = "engineering"
      # type    = "radio"
      type = "checkbox"

      content {
        value {
          option {
            label = "engineering namespace #1"
            value = "eng-ns-1"
          }
          option {
            label   = "engineering namespace #2"
            value   = "eng-ns-2"
            checked = true
          }
        }
      }
    }

    conditional {
      source    = "kms.name"
      condition = "finances"
      type      = "radio"

      content {
        value {
          option {
            label = "finances namespace #1"
            value = "fin-ns-1"
          }
          option {
            label   = "finances namespace #2"
            value   = "fin-ns-2"
            checked = true
          }
        }
      }
    }

    conditional {
      source    = "kms.name"
      condition = "sales"
      content {
        static = "sales-ns-1"
      }
    }
  }
}

resource "autocloud_blueprint" "example" {
  name = "KMS conditionals"

  ###
  # UI Configuration
  #
  author       = "enrique.enciso@autocloud.dev"
  description  = "Terraform Generator storage in cloud"
  instructions = <<-EOT
  To deploy this generator, follow these simple steps:

  step 1: step-1-description
  step 2: step-2-description
  step 3: step-3-description
  EOT

  labels = ["aws"]

  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new EKS generator, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new EKS generator, created by {{authorName}}"
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
    destination = "kms.tf"
    variables = {
    }

    modules = [
      autocloud_module.kms.name ## change for id? instead of name
    ]
  }
  config = data.autocloud_blueprint_config.kms_custom_form.config
}
