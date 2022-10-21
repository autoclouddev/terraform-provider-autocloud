package autocloud_provider

import (
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
  #
  #file {
  #  action = "CREATE"
  #
  #  path_from_root = "some-path"
  #
  #  filename_template = "blueprint-file-{{clusterName}}.tf"
  #	config = [{
  #		module = autocloud_module.s3_bucket.id
  #		filename_vars = {
  #			clusterName = "bucketName"
  #		}
  #
  #	}]
  #
  #}


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
      body                    = ""
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
				),
			},
		},
	})
}
