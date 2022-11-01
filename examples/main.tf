
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
provider "autocloud" {}

data "autocloud_me" "current_user" {}

data "autocloud_github_repos" "repos" {}
locals {
  ####
  # Destination repos
  dest_repos = data.autocloud_github_repos.repos.data[*].url
}



resource "autocloud_module" "s3_bucket" {

  ####
  # Name of the generator
  name = "S3Bucket"

  ####
  # Can be any supported terraform source reference, must optionaly take version
  #
  #   source = "app.terraform.io/autocloud/aws/s3_bucket"
  #   version = "0.24.0"
  #
  # See docs: https://developer.hashicorp.com/terraform/language/modules/sources

  version = "3.4.0"
  source  = "terraform-aws-modules/s3-bucket/aws"

}



resource "autocloud_module" "cloudfront" {

  ####
  # Name of the generator
  name = "Cloudfront"

  ####
  # Can be any supported terraform source reference, must optionaly take version
  #
  #   source = "app.terraform.io/autocloud/aws/s3_bucket"
  #   version = "0.24.0"
  #
  # See docs: https://developer.hashicorp.com/terraform/language/modules/sources

  version = "3.0.0"
  source  = "terraform-aws-modules/cloudfront/aws"

}

resource "autocloud_blueprint" "example" {
  name = "S3andEKS"

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
    destination_branch = "master"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new EKS generator, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new EKS generator, created by {{authorName}}"
      body                    = jsonencode(file("./files/pull_request.md.tpl"))
      variables = {
        authorName  = "generic.authorName"
        clusterName = "S3Bucket.Bucket"
        dummyParam = autocloud_module.s3_bucket.variables["restrict_public_buckets"]
      }
    }
  }


  ###
  # File definitions
  #
  file {
    action = "CREATE"

    path_from_root = ""

    filename_template = "s3-bucket-{{Bucket}}.tf"
    filename_vars = {
      Bucket = "S3Bucket.Bucket"
    }
    modules = ["S3Bucket", "Cloudfront"]
  }


  autocloud_module {
    id = autocloud_module.s3_bucket.id

    form_config     = templatefile("${path.module}/files/s3bucket.vars.tpl", {})  # example from file
    # form_config     = autocloud_module.s3_bucket.form_config                    # example from resource
    # form_config     = data.autocloud_form_config.s3_bucket.form_config          # example from data

    template_config = file("${path.module}/files/s3bucket.tf.tpl")
  }

  autocloud_module {
    id = autocloud_module.cloudfront.id
  }
}
