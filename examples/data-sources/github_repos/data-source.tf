terraform {
  required_providers {
    autocloud = {
      source  = "autoclouddev/autocloud"
      version = "~>0.9"
    }
  }
}


data "autocloud_github_repos" "example" {}


###
# Will contain a list of repo objects similar to this:
#
# [
#   {
#       id = "611893545"
#       name = "terraform-provider-autocloud"
#       url = "https://github.com/autoclouddev/terraform-provider-autocloud"
#       description = "Terraform Provider Repo Example"
#   }
# ]

output "repos" {
  description = "List of git repositories AutoCloud can submit pull requests to."
  value = data.autocloud_github_repos.example.data
}
