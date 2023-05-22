terraform {
  required_providers {
    autocloud = {
      version = "0.2.0"
      source  = "autocloud.io/autoclouddev/autocloud"
    }
  }
}

data "autocloud_github_repos" "example" {
    data = {
        id = "611893545"
        name = "terraform-provider-autocloud"
        url = "https://github.com/autoclouddev/terraform-provider-autocloud"
        description = "Terraform Provider Repo Example"
    }
}