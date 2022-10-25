
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
provider "autocloud" {}

# module "test" {
#   source = "./autocloud"
# }

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


# module "cloud-storage" {
#   source  = "terraform-google-modules/cloud-storage/google"
#   version = "3.4.0"
#   # insert the 3 required variables here
# }

resource "autocloud_module_1" "s3_bucket" {

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
	source = "terraform-aws-modules/s3-bucket/aws"

}

resource "autocloud_module_1" "eks" {

	####
	# Name of the generator
	name = "EKS"

	####
	# Can be any supported terraform source reference, must optionaly take version
	#
	#   source = "app.terraform.io/autocloud/aws/s3_bucket"
	#   version = "0.24.0"
	#
	# See docs: https://developer.hashicorp.com/terraform/language/modules/sources

	version = "2.0.2"
	source = "howdio/eks/aws"

}

# module "s3-bucket" {
#   source  = "terraform-aws-modules/s3-bucket/aws"
#   version = "3.4.0"
# }
resource "autocloud_module" "example" {
  name = "example_s3"

  ###
  # UI Configuration
  #
  author       = "enrique.enciso@autocloud.dev"
  slug         = "example_s3"
  description  = "Terraform Generator storage in cloud"
  instructions = <<-EOT
  To deploy this generator, follow these simple steps:

  step 1: step-1-description
  step 2: step-2-description
  step 3: step-3-description
  EOT

  labels = ["aws"]

  ###
  # TF source
  #
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "3.4.0"
  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "master"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new EKS generator, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new EKS generator, created by {{authorName}}"
      body                    = jsonencode(file("./generator/pull_request.md.tpl"))
      variables = {
        authorName  = "generic.authorName"
        clusterName = "ExampleS3.Bucket"
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
      Bucket = "ExampleS3.Bucket"
    }
  }

  generator_config_location = "local"
  generator_config_json     = templatefile("${path.module}/form.autocloud.json.tpl", {})

  autocloud_module_1 {
    id = autocloud_module_1.s3_bucket.id
  }

  autocloud_module_1 {
    id = autocloud_module_1.eks.id
  }
}

# output "test" {
#   value = module.test.autocloud_me_output
# }

output "terraform_template" {
  value = autocloud_module.example.template
}


# output "repos" {
#   value = module.test.autocloud_github_repos
# }
