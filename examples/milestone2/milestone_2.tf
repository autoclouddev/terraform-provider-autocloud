
##
# repos previously connected to autocloud
data "autocloud_github_repos" "repos" {}


##
# Getting a previous uploaded module
data "autocloud_module" "cloudfront_distribution" {
  version = "1.0.2" ## a user could upload multiple versions, we are fetching this version
}



#uploading a new module
resource "autocloud_module" "s3_bucket" {
  ###
  # Name of the generator
  name = "s3_bucket"



  ###
  # Can be any supported terraform source reference, must optionaly take version
  #
  #   source = "app.terraform.io/autocloud/aws/compute/eks/control_plane"
  #   version = "0.24.0"
  #

  source  = "some_private_registry.towhich.we_have_access_to.com/modules/s3" // paste a url of a public accesible module from a registry
  version = "1.0.1"


}



##
# this helps in redenring the form, takes the default form and template calculated when the user uploads the module
# it has two outputs
# usage_template, wich is the output of the generated code
# form_config, wich is the UI form

module "form_builder_s3" {
  source  = "tfe.autocloud.dev/utlis/formbuilder" #this module only makes computations locally, this is public
  version = "1.0.0"
  form    = autocloud_module.s3_bucket.form_config
  # distinct will keep the first duplicated element
  # here we are overwritting the variable "bucket_name"
  variables = distinct(concat([
    {
      id   = autocloud_module.s3_bucket.variables["bucket_name"].id
      type = "radio",
      options = [
        {
          value = "nonprod"
          label = "non prod account"
        },
        {
          value = "prod"
          label = "prod account"
        }
      ]
    },

    ],
    autocloud_module.s3_bucket.variables
    )
  )
}



resource "autocloud_blueprint" "static_site_generator" {
  ###
  # UI Configuration
  #
  name         = "static_site_generator"
  author       = "enrique.enciso@autocloud.dev"
  display_name = "(AutoCloud) EKS generator"
  #slug         = "autocloud_eks_generator" // this will be processed in the sdk, this should not exists in the schema
  description  = "Terraform Generator for a static site"
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
  # Modules
  #

  autocloud_module {
    id              = autocloud_module.s3_bucket.id
    template_config = autocloud_module.s3_bucket.template
    form_config     = autocloud_module.s3_bucket.form_config
  }

  autocloud_module {
    id              = data.autocloud_module.cloudfront_distribution.id
    template_config = data.autocloud_module.cloudfront_distribution.template_config
    form_config     = data.autocloud_module.cloudfront_distribution.form_config
  }


  ###
  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = data.autocloud_github_repos.repos
    git_url_default = local.dest_repos[0] # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new static site {{siteName}} , created by {{authorName}}"
      commit_message_template = "[AutoCloud] new static site, created by {{authorName}}"
      body                    = jsonencode(file("./milestone_1_pull_request.md.tpl")) // not sure if we should use jsonencode
      variables = {
        authorName = "generic.authorName",
        siteName   = autocloud_module.s3_bucket.variables["bucket_name"].id
      }
    }
  }



  ###
  # File definitions
  #
  file {
    action = "CREATE" // if this is not defined, lets place "CREATE" as default

    path_from_root = "eks/" // assumes there is folder called eks, add a validation for a later milestone

    filename_template = "static-site-{{siteName}}.tf"
    filename_vars = {
      siteName = autocloud_module.s3_bucket.variables["bucket_name"].id
    }
  }


}
