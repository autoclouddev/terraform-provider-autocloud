package autocloud_provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccAutocloudBluePrint = `
resource "autocloud_blueprint" "bar" {
  name = "FirstBluePrint"
  author = "enrique.enciso@autocloud.dev"
  description  = "Terraform Generator for Elastic Kubernetes Service"
  instructions = <<-EOT
  To deploy this generator, follow these simple steps:

  step 1: step-1-description
  step 2: step-2-description
  step 3: step-3-description
  EOT
  labels       = [
	"aws"
  ]

  ###
  # File definitions
  # THIS HAS TO CHANGE TO SUPPORT FILE PER MODULE
  #
  file {
    action = "CREATE"

    path_from_root = "some-path"

    filename_template = "eks-cluster-{{clusterName}}.tf"
    filename_vars = {
      clusterName = "EKSGenerator.clusterName"
    }
  }


  ###
  # Destination repository git configuration
  #
  git_config {
    destination_branch = "main"

    git_url_options = ["github.com/autoclouddev/terraform-generator-test"]
    git_url_default = "github.com/autoclouddev/terraform-generator-test"

    pull_request {
      title                   = "[AutoCloud] new static site {{siteName}} , created by {{authorName}}"
      commit_message_template = "[AutoCloud] new static site, created by {{authorName}}"
      body                    = "Body Example"
      variables = {
        authorName = "generic.authorName",
        siteName   = "generic.SiteName"  #autocloud_module.s3_bucket.variables["bucket_name"].id
      }
    }
  }


  ###
  # Modules
  #
  #autocloud_module {
  #  id              = autocloud_module.s3_bucket.id
  #  template_config = autocloud_module.s3_bucket.template
  #  form_config     = data.autocloud_form_config.s3_bucket
  #}


}
`

func TestAccAutocloudBlueprint(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAutocloudBluePrint,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "name", "FirstBluePrint"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "author", "enrique.enciso@autocloud.dev"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "description", "Terraform Generator for Elastic Kubernetes Service"),
					resource.TestMatchResourceAttr(
						"autocloud_blueprint.bar", "instructions", regexp.MustCompile("To deploy this generator, follow these simple steps")),
					resource.TestCheckResourceAttrSet(
						"autocloud_blueprint.bar", "labels.0"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "file.0.action", "CREATE"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "file.0.path_from_root", "some-path"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "file.0.filename_template", "eks-cluster-{{clusterName}}.tf"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "file.0.filename_vars.clusterName", "EKSGenerator.clusterName"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "git_config.0.destination_branch", "main"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "git_config.0.git_url_default", "github.com/autoclouddev/terraform-generator-test"),
					resource.TestCheckResourceAttrSet(
						"autocloud_blueprint.bar", "git_config.0.git_url_options.0"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "git_config.0.pull_request.0.title", "[AutoCloud] new static site {{siteName}} , created by {{authorName}}"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "git_config.0.pull_request.0.commit_message_template", "[AutoCloud] new static site, created by {{authorName}}"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "git_config.0.pull_request.0.body", "Body Example"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "git_config.0.pull_request.0.variables.authorName", "generic.authorName"),
					resource.TestCheckResourceAttr(
						"autocloud_blueprint.bar", "git_config.0.pull_request.0.variables.siteName", "generic.SiteName"),
				),
			},
		},
	})
}
