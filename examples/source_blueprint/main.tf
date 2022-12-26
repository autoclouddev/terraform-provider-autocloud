terraform {
  required_providers {
    autocloud = {
      version = "0.2.0"
      source  = "autocloud.io/autoclouddev/autocloud"
    }
  }
}


/*
USE THIS FILE AS YOU NEED FIT, THIS IS JUST A PLAYGROUND

*/
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
    if length(regexall("/infrastructure-live-demo", repo)) > 0
  ]
}

resource "autocloud_module" "kms" {
  name   = "kms"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/security/kms?ref=0.1.0"
}


data "autocloud_blueprint_config" "generic" {
  variable {
    name         = "env"
    display_name = "environment target"
    helper_text  = "environment target description"
    form_config {
      type = "radio"
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
}

data "autocloud_blueprint_config" "kms_custom_form" {
  source = {
    kms     = autocloud_module.kms.blueprint_config
    generic = data.autocloud_blueprint_config.generic.blueprint_config
  }

  variable {
    name         = "description"
    display_name = "choose the usage of this kms"
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
      }
    }
  }
  //add override to test backward comp

  ### overriding "namespace" (string)
  # hardcoded value
  variable {
    name  = "namespace"
    value = "some-namespace"
  }

  ### overriding "enable_key_rotation" (bool)
  # hardcoded value
  # variable {
  #   name  = "enable_key_rotation"
  #   value = false
  # }

  # form override
  variable {
    name         = "enable_key_rotation"
    display_name = "enable_key_rotation display name"
    helper_text  = "enable_key_rotation helper text"
    form_config {
      type = "radio"
      options {
        option {
          label   = "ENABLE"
          value   = true
          checked = false
        }
        option {
          label   = "DISABLE"
          value   = false
          checked = true
        }
      }
    }
  }

  ### overriding "deletion_window_in_days" (number)
  # hardcoded value
  # variable {
  #   name    = "deletion_window_in_days"
  #   value   = 30
  # }

  # form override
  variable {
    name         = "deletion_window_in_days"
    display_name = "deletion_window_in_days display name"
    helper_text  = "deletion_window_in_days helper text"
    form_config {
      type = "radio"
      options {
        option {
          label   = "30 days"
          value   = 30
          checked = false
        }
        option {
          label   = "60 days"
          value   = 60
          checked = false
        }
        option {
          label   = "90 days"
          value   = 90
          checked = true
        }
      }
    }
  }
}




resource "autocloud_blueprint" "example" {
  name = "KMS from examples"

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
        authorName = "generic.authorName"
        //clusterName = "resource.autocloud_blueprint_config.kms.variables["name"]"
        clusterName = "kms.name"
        //dummyParam  = autocloud_module.s3_bucket.variables["restrict_public_buckets"]
      }
    }
  }


  ###
  # File definitions
  #

  file {
    action      = "CREATE"
    destination = "{{kmsName}}.tf"
    variables = {
      kmsName = "kms.name" # variables["name"]
    }

    modules = [
      autocloud_module.kms.name ## change for id? instead of name
    ]
  }
  config = data.autocloud_blueprint_config.kms_custom_form.config
}
