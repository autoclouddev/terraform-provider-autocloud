package autocloud_provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAutocloudModule(t *testing.T) {
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
					resource.TestCheckResourceAttr(
						"autocloud_module.foo", "instructions", "Instructions text"),
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
  instructions = "Instructions text"
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

}
`
