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
provider "autocloud" {
  # endpoint = "https://api.autocloud.domain.com/api/v.0.0.1"
}

//data "autocloud_github_repos" "repos" {}


data "autocloud_dummy" "foobar" {
  name = "name"
  values = {
    hello = "world"
  }
  #values = "hello"
  //values = 4
}

output "state_output" {
  value = data.autocloud_dummy.foobar.values
}
