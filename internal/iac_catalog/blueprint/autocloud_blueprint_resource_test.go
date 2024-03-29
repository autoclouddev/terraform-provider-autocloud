package blueprint_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/acctest"
)

const testAccAutocloudBluePrint = `

data "autocloud_blueprint_config" "test" {
	source = {}

	omit_variables = []

	variable {
	  name         = "app_name"
	  display_name = "Application Name"
	  type         = "shortText"

	  validation_rule {
		rule          = "isRequired"
		error_message = "You must provide an application name to provision a resource"
	  }
	}
  }

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

    destination = "eks-cluster-{{clusterName}}.tf"
    variables = {
      clusterName = "EKSGenerator.clusterName"
    }
	modules = ["EKSGenerator"]
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
  config = data.autocloud_blueprint_config.test.config

}
`

func TestAccAutocloudBlueprint(t *testing.T) {
	//t.SkipNow()
	resourceName := "autocloud_blueprint.bar"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAutocloudBluePrint,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", "FirstBluePrint"),
					resource.TestCheckResourceAttr(
						resourceName, "author", "enrique.enciso@autocloud.dev"),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Terraform Generator for Elastic Kubernetes Service"),
					resource.TestMatchResourceAttr(
						resourceName, "instructions", regexp.MustCompile("To deploy this generator, follow these simple steps")),
					resource.TestCheckResourceAttrSet(
						resourceName, "labels.0"),
					resource.TestCheckResourceAttr(
						resourceName, "file.0.action", "CREATE"),
					resource.TestCheckResourceAttr(
						resourceName, "file.0.destination", "eks-cluster-{{clusterName}}.tf"),
					resource.TestCheckResourceAttr(
						resourceName, "file.0.variables.clusterName", "EKSGenerator.clusterName"),
					resource.TestCheckResourceAttr(
						resourceName, "file.0.modules.0", "EKSGenerator"),
					resource.TestCheckResourceAttr(
						resourceName, "git_config.0.destination_branch", "main"),
					resource.TestCheckResourceAttr(
						resourceName, "git_config.0.git_url_default", "github.com/autoclouddev/terraform-generator-test"),
					resource.TestCheckResourceAttrSet(
						resourceName, "git_config.0.git_url_options.0"),
					resource.TestCheckResourceAttr(
						resourceName, "git_config.0.pull_request.0.title", "[AutoCloud] new static site {{siteName}} , created by {{authorName}}"),
					resource.TestCheckResourceAttr(
						resourceName, "git_config.0.pull_request.0.commit_message_template", "[AutoCloud] new static site, created by {{authorName}}"),
					resource.TestCheckResourceAttr(
						resourceName, "git_config.0.pull_request.0.body", "Body Example"),
					resource.TestCheckResourceAttr(
						resourceName, "git_config.0.pull_request.0.variables.authorName", "generic.authorName"),
					resource.TestCheckResourceAttr(
						resourceName, "git_config.0.pull_request.0.variables.siteName", "generic.SiteName"),
					resource.TestCheckResourceAttrSet(
						resourceName, "config"),
				),
			},
		},
	})
}

func TestAutocloudBlueprintHasAtMostOneGitConfigError(t *testing.T) {
	expectedError := `No more than 1 "git_config" blocks are allowed`
	terraform := `resource "autocloud_blueprint" "bar" {
		git_config {}
		git_config {}
	  }`
	acctest.ValidateErrors(t, fmt.Errorf(expectedError), terraform)
}

func TestAutocloudBlueprintHasAtMostOneFileBlockError(t *testing.T) {
	expectedError := `Insufficient file blocks`
	terraform := `resource "autocloud_blueprint" "bar" {
		git_config {}
	  }`
	acctest.ValidateErrors(t, fmt.Errorf(expectedError), terraform)
}

func TestAutocloudBlueprinntHasMissingAttributesOnFileBlockError(t *testing.T) {
	expectedError := `file block should contain content or modules attributes`
	terraform := `resource "autocloud_blueprint" "bar" {
		file {
			action      = "CREATE"
			destination = "s3.tf"
			variables = {}
		}
	  }`
	acctest.ValidateErrors(t, fmt.Errorf(expectedError), terraform)
}

func TestAutocloudBlueprinntHasMissingModulesWithHeader(t *testing.T) {
	expectedError := `modules can not be empty when using header or footer attributes`
	terraform := `resource "autocloud_blueprint" "bar" {
		file {
			action      = "CREATE"
			destination = "s3.tf"
			variables = {}
			header = "# header content"
		}
	  }`
	acctest.ValidateErrors(t, fmt.Errorf(expectedError), terraform)
}

func TestAutocloudBlueprinntHasMissingModulesWitFooter(t *testing.T) {
	expectedError := `modules can not be empty when using header or footer attributes`
	terraform := `resource "autocloud_blueprint" "bar" {
		file {
			action      = "CREATE"
			destination = "s3.tf"
			variables = {}
			footer = "# footer content"
		}
	  }`
	acctest.ValidateErrors(t, fmt.Errorf(expectedError), terraform)
}

func minimal_blueprint() string {
	return `
	data "autocloud_blueprint_config" "test" {
		source = {}

		omit_variables = []

		variable {
		  name         = "app_name"
		  display_name = "Application Name"
		 type         = "shortText"

		  validation_rule {
			rule          = "isRequired"
			error_message = "You must provide an application name to provision a resource"
		  }
		}
	  }

	resource "autocloud_blueprint" "minimal" {
	  name = "FirstBluePrint"
	  author = "enrique.enciso@autocloud.dev"

	  ###
	  # File definitions
	  # THIS HAS TO CHANGE TO SUPPORT FILE PER MODULE
	  #
	  file {
		action = "CREATE"

		destination = "eks-cluster-{{clusterName}}.tf"
		variables = {
		  clusterName = "EKSGenerator.clusterName"
		}
		#modules = ["EKSGenerator"]
		content ="hellO"
	  }

	  config = data.autocloud_blueprint_config.test.config

	}

	`
}

func TestAccAutocloudMinimalBlueprint(t *testing.T) {
	//t.SkipNow()
	resourceName := "autocloud_blueprint.minimal"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: minimal_blueprint(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						resourceName, "name"),
					resource.TestCheckResourceAttrSet(
						resourceName, "author"),
					resource.TestCheckResourceAttrSet(
						resourceName, "config"),
					resource.TestCheckNoResourceAttr(resourceName, "git_config.#"),
				),
			},
		},
	})
}
