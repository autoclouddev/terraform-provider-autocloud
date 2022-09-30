
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

module "test" {
  source = "./autocloud"
}

# uncomment this to test milestone1
# module "milestone_1" {
#   source = "./milestone1"
# }



# module "cloud-storage" {
#   source  = "terraform-google-modules/cloud-storage/google"
#   version = "3.4.0"
#   # insert the 3 required variables here
# }


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

}

output "test" {
  value = module.test.autocloud_me_output
}

output "terraform_template" {
  value = autocloud_module.example.template
}


output "repos" {
  value = module.test.autocloud_github_repos
}
