terraform {
  required_providers {
    autocloud = {
      version = "0.2"
      source  = "autocloud.io/autocloud/autocloud"
    }
  }
}



data "autocloud_me" "current_user" {}

data "autocloud_github_repos" "repos" {}



# Only returns email
output "autocloud_me_output" {
  value = {
    email : data.autocloud_me.current_user.email
  }

}

output "autocloud_github_repos" {
  value = {
    repos : "${length(data.autocloud_github_repos.repos)}"
  }

}
