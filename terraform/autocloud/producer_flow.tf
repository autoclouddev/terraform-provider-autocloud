
/*
This file is the final version of the producer flow
It is commented because it is not supported yet
*/


/*
locals {
  ####
  # VPC IDs of VPCs to allow in form question
  vpc_ids = [
    "vpc-04ba69a71c19bcd51",
    "vpc-04ba69a71c19bcd51"
  ]


  ####
  # Destination repos
  dest_repos = [
    for repo in data.autocloud_github_integrations.repos :
    # Filter out the repository URLs that contain "-live" in the URL
    # This  is a ridiculous way to determine if a string contains a substring, but HCL is a silly language
    repo if length(regexall(".*-live.*", repo.git_url)) > 0
  ]
}





####
# Github resources
#

# Would be a good way to define Github integrations
resource "autocloud_github_integration" "src_iac_code" {
  git_url      = "git@github.com:autoclouddev/self-hosted-infrastructure-modules.git"
  access_token = var.github_token
  display_name = "Source Modules"
}

resource "autocloud_github_integration" "saas_dest_iac_code" {
  git_url      = "git@github.com:autoclouddev/self-hosted-infrastructure-live.git"
  access_token = var.github_token
  display_name = "Live Infrastructure"
}

resource "autocloud_github_integration" "self_hossted_dest_iac_code" {
  git_url      = "git@github.com:autoclouddev/infrastructure-live.git"
  access_token = var.github_token
  display_name = "Live Infrastructure"
}

####
# Remote module w/ generator config in code
#
resource "autocloud_module" "eks_generator_remote" {
  ###
  # Name of the generator
  name = "EKSGenerator"

  ###
  # Can be any supported terraform source reference, must optionaly take version
  # 
  #   source = "app.terraform.io/autocloud/aws/compute/eks/control_plane"
  #   version = "0.24.0"
  #
  source = "${autocloud_github_integration.src_iac_code.git_url}//aws/compute/eks/control_plane?ref=0.24.0"

  ###
  # UI Configuration
  #
  author       = "enrique.enciso@autocloud.dev"
  display_name = "(AutoCloud) EKS generator"
  slug         = "autocloud_eks_generator"
  description  = "Terraform Generator for Elastic Kubernetes Service"
  instructions = <<-EOT
  To deploy this generator, follow these simple steps:

  step 1: collect underpants
  step 2: ????
  step 3: profit
  EOT
  providers = [
    "aws"
  ]

  ###
  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = local.dest_repos
    git_url_default = local.dest_repos[0] # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new EKS generator, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new EKS generator, created by {{authorName}}"
      body                    = jsonencode(file("./generator/pull_request.md.tpl"))
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
    action = "CREATE"

    path_from_room = ""

    filename_template = "eks-cluster-{{clusterName}}.tf"
    filename_vars = {
      clusterName = "EKSGenerator.clusterName"
    }
  }

  generator_config_location = "module"
  generator_config_path     = "path/in/repo/to/config.json"
}





####
# Remote module w/ generator config specified locally
#

data "aws_vpc" "options" {
  for_each = toset(locals.vpc_ids)
  id       = each.value
}

data "aws_subnets" "options" {
  for_each = toset(locals.vpc_ids)

  filter {
    name   = "vpc-id"
    values = [each.value]
  }
}

data "template_file" "eks_generator_config" {
  template = file("./generator/eks_generator.autocloud.json.tpl")
  vars = {
    vpc_field_options = jsonencode([
      for vpc in data.aws_vpc.options :
      {
        label   = vpc.tags["Name"]
        fieldId = vpc.tags["Name"]
        value   = vpc.id
        checked = false
        dependencyData = [
          {
            dependentId = "EKSGenerator.vpcSubnets"
            type        = "fieldOptions"
            values = [
              for subnet in data.aws_subnets.options[vpc.id] :
              {
                "label" : subnet.id
                "fieldId" : subnet.id
                "value" : subnet.id
                "checked" : false
              }
            ]
          }
        ]
      }
    ])
  }
}

resource "autocloud_module" "eks_generator_local" {
  ###
  # Name of the generator
  name = "EKSGenerator"

  ###
  # Same source options as above
  #
  source = "${autocloud_github_integration.src_iac_code.git_url}//aws/compute/eks/control_plane?ref=0.24.0"


  ###
  # UI Configuration
  #
  author       = "enrique.enciso@autocloud.dev"
  display_name = "(AutoCloud) EKS generator"
  slug         = "autocloud_eks_generator"
  description  = "Terraform Generator for Elastic Kubernetes Service"
  instructions = <<-EOT
  To deploy this generator, follow these simple steps:

  step 1: collect underpants
  step 2: ????
  step 3: profit
  EOT
  providers = [
    "aws"
  ]

  ###
  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = local.dest_repos
    git_url_default = local.dest_repos[0] # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new EKS generator, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new EKS generator, created by {{authorName}}"
      body                    = jsonencode(file("./generator/pull_request.md.tpl"))
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
    action = "CREATE"

    path_from_room = ""

    filename_template = "eks-cluster-{{clusterName}}.tf"
    filename_vars = {
      clusterName = "EKSGenerator.clusterName"
    }
  }

  generator_config_location = "local"
  generator_config_json     = data.template_file.eks_generator.rendered
}

*/




