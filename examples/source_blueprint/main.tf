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

resource "autocloud_module" "kms" {
  name   = "kms"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/security/kms?ref=0.1.0"
}

data "autocloud_blueprint_config" "kms_custom_form" {
  source = {
    kms = autocloud_module.kms.blueprint_config
    s3  = autocloud_module.kms.blueprint_config
  }

  variable {
    name = "source.kms.variables.key_name"
    conditional {
      source   = "source.s3.variables.name" # reference syntax
      conditon = "prod"

      content {
        value = "hello"
      }
    }
  }

  //add override to test backward comp
}



output "form_module" {
  value = autocloud_module.kms.blueprint_config
}

output "form_blueprint" {
  value = data.autocloud_blueprint_config.kms_custom_form.blueprint_config
}
/*
resource "autocloud_blueprint" "example" {
  name = "S3andCloudfront"

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
      body                    = file("./files/pull_request.md.tpl")
      variables = {
        authorName  = "generic.authorName"
        clusterName = resource.autocloud_blueprint_config.kms.variables["name"]
        clusterName = "S3Bucket.Bucket"
        dummyParam  = autocloud_module.s3_bucket.variables["restrict_public_buckets"]
      }
    }
  }


  ###
  # File definitions
  #

  file {
    action    = "CREATE"
    dest_path = "{{kmsName}}.tf"
    variables = {
      kmsName = data.autocloud_form_config.kms_custom_form.variables.name # variables["name"]
    }

    modules = [
      autocloud_module.kms.name ## change for id? instead of name
    ]
  }
  config = data.autocloud_blueprint_config.kms_custom_form.blueprint_config
}



*/
