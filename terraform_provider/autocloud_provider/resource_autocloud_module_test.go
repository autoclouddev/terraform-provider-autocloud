package autocloud_provider

import (
	"regexp"
	"testing"

	"github.com/joho/godotenv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAutocloudModule(t *testing.T) {
	godotenv.Load("../../.env")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAutocloudModule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"autocloud_module.foo", "name", "EKSGenerator"),
					resource.TestCheckResourceAttr(
						"autocloud_module.foo", "author", "enrique.enciso@autocloud.dev"),
					resource.TestCheckResourceAttr(
						"autocloud_module.foo", "slug", "autocloud_eks_generator"),
					resource.TestCheckResourceAttr(
						"autocloud_module.foo", "description", "Terraform Generator for Elastic Kubernetes Service"),
					resource.TestMatchResourceAttr(
						"autocloud_module.foo", "instructions", regexp.MustCompile(`(.\s)*To deploy this generator, follow these simple steps.*`)),
					resource.TestCheckResourceAttr(
						"autocloud_module.foo", "generator_config_location", "local"),
					resource.TestMatchResourceAttr(
						"autocloud_module.foo", "generator_config_json", regexp.MustCompile(".*formQuestion.*")),
					resource.TestCheckTypeSetElemAttr(
						"autocloud_module.foo", "labels.*", "aws"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"autocloud_module.foo", "file.*", map[string]string{
							"action":            "CREATE",
							"path_from_root":    "some-path",
							"filename_template": "eks-cluster-{{clusterName}}.tf",
						}),
					resource.TestCheckResourceAttr(
						"autocloud_module.foo", "file.0.filename_vars.clusterName", "EKSGenerator.clusterName"),
				),
			},
		},
	})
}

const testAccAutocloudModule = `
resource "autocloud_module" "foo" {
  name 		   = "EKSGenerator"
  author       = "enrique.enciso@autocloud.dev"
  slug         = "autocloud_eks_generator"
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
  file {
    action = "CREATE"

    path_from_root = "some-path"

    filename_template = "eks-cluster-{{clusterName}}.tf"
    filename_vars = {
      clusterName = "EKSGenerator.clusterName"
    }
  }

  generator_config_location = "local"
  generator_config_json     = <<-EOT
  {
	"terraformModules": {
	  "EKSGenerator": [
		{
		  "id": "EKSGenerator.clusterName",
		  "module": "EKSGenerator",
		  "type": "string",
		  "formQuestion": {
			"fieldId": "EKSGenerator.clusterName",
			"fieldType": "shortText",
			"fieldLabel": "Cluster name",
			"validationRules": [
			  {
				"errorMessage": "This field is required",
				"rule": "isRequired"
			  }
			]
		  }
		}
	  ]
	}
  }
  EOT
}
`
