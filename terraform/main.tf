
terraform {
  required_providers {
    autocloud = {
      version = "0.2"
      source  = "autocloud.io/autocloud/autocloud"
    }
  }
}
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

output "test" {
  value = module.test.autocloud_me_output
}
