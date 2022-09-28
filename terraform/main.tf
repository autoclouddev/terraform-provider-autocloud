
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
  username = ""
  password = ""
}

module "test" {
  source = "./autocloud"
}

# uncomment this to test milestone1
# module "milestone_1" {
#   source = "./milestone1"
# }


resource "autocloud_module" "example" {
  name = "example"

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

}

output "test" {
  value = module.test.autocloud_me_output
}
