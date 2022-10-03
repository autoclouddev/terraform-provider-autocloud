
terraform {
  required_providers {
    autocloud = {
      version = "0.2.0"
      source  = "autocloud.io/autocloud/autocloud"
    }
  }
}


/*
USE THIS FILE AS YOU NEED FIT, THIS IS JUST A PLAYGROUND

*/
provider "autocloud" {
  username = "enrique.enciso@autocloud.dev"
  password = "publisherFlow1#"
}

module "test" {
  source = "./autocloud"
}

# uncomment this to test milestone1
# module "milestone_1" {
#   source = "./milestone1"
# }

data "autocloud_me" "current_user" {}

data "autocloud_github_repos" "repos" {}
locals {
  ####
  # Destination repos
  dest_repos = data.autocloud_github_repos.repos.data[*].url
}

resource "autocloud_module" "example" {
  name = "example"
  module_name = "EKSGenerator"

  ###
  # UI Configuration
  #
  author       = "enrique.enciso@autocloud.dev"
  slug         = "autocloud_eks_generator"
  description  = "Terraform Generator for Elastic Kubernetes Service"
  instructions = <<-EOT
  To deploy this generator, follow these simple steps:

  step 1: step-1-description
  step 2: step-2-description
  step 3: step-3-description
  EOT

  labels = ["aws"]

  ###
  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "master"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0]  : "" # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new EKS generator, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new EKS generator, created by {{authorName}}"
      body                    = jsonencode(file("./generator/pull_request.md.tpl"))
      variables = {
        authorName  = "generic.authorName"
        clusterName = "EKSGenerator.clusterName"
      }
    }
  }


  ###
  # File definitions
  #
  file {
    action = "CREATE"

    path_from_root = ""

    filename_template = "eks-cluster-{{clusterName}}.tf"
    filename_vars = {
      clusterName = "EKSGenerator.clusterName"
    }
  }

}

output "test" {
  value = module.test.autocloud_me_output
}

output "repos" {
  value = module.test.autocloud_github_repos
}
