
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

output "test" {
  value = module.test.autocloud_me_output
}

output "repos" {
  value = module.test.autocloud_github_repos
}
