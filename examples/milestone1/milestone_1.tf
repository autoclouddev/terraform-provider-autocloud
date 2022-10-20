data "autocloud_github_repos" "repos" {}


resource "autocloud_module" "eks_generator_local" {
  ###
  # Name of the generator
  name = "EKSGenerator"



  ###
  # Can be any supported terraform source reference, must optionaly take version
  #
  #   source = "app.terraform.io/autocloud/aws/compute/eks/control_plane"
  #   version = "0.24.0"
  #

  source  = "" // paste a url of a public accesible module from a registry
  version = ""


  ###
  # UI Configuration
  #
  author       = "enrique.enciso@autocloud.dev"
  display_name = "(AutoCloud) EKS generator"
  #slug         = "autocloud_eks_generator" // this will be processed in the sdk, this should not exists in the schema
  description  = "Terraform Generator for Elastic Kubernetes Service"
  instructions = <<-EOT
  To deploy this generator, follow these simple steps:

  step 1: collect underpants
  step 2: ????
  step 3: profit
  EOT
  labels = [
    "aws"
  ]



  ###
  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = data.autocloud_github_repos.repos
    git_url_default = local.dest_repos[0] # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new EKS generator, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new EKS generator, created by {{authorName}}"
      body                    = jsonencode(file("./milestone_1_pull_request.md.tpl")) // not sure if we should use jsonencode
      variables = {
        authorName  = "generic.authorName",
        clusterName = "EKSGenerator.clusterName"
      }
    }
  }



  ###
  # File definitions
  #
  file {
    action = "CREATE" // if this is not defined, lets place "CREATE" as default

    path_from_root = "eks/" // assumes there is folder called eks, add a validation for a later milestone

    filename_template = "eks-cluster-{{clusterName}}.tf"
    filename_vars = {
      clusterName = "EKSGenerator.clusterName"
    }
  }

  #FORM
  generator_config_location = "local"
  generator_config_path     = "./form_config.json"
}

