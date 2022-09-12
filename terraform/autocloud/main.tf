terraform {
  required_providers {
    autocloud = {
      version = "0.2"
      source  = "autocloud.io/autocloud/autocloud"
    }
  }
}



data "autocloud_me" "current_user" {}



# Only returns email
output "autocloud_me_output" {
  value = {
    email : data.autocloud_me.current_user.email
  }

}
