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
						"autocloud_module.foo", "name", "(AutoCloud) EKS generator"),
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
				),
			},
		},
	})
}

const testAccAutocloudModule = `
resource "autocloud_module" "foo" {
  name 		   = "(AutoCloud) EKS generator"
  author       = "enrique.enciso@autocloud.dev"
  slug         = "autocloud_eks_generator"
  description  = "Terraform Generator for Elastic Kubernetes Service"
  instructions = "Instructions text"
  labels    = [ 
	"aws"
  ]  
}
`
