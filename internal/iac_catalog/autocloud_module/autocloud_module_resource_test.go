package autocloud_module_test

import (
	"testing"

	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccAutocloudModule = `
resource "autocloud_module" "s3_bucket" {

	####
	# Name of the generator
	name = "S3Bucket"

	####
	# Can be any supported terraform source reference, must optionaly take version
	#
	#   source = "app.terraform.io/autocloud/aws/s3_bucket"
	#   version = "0.24.0"
	#
	# See docs: https://developer.hashicorp.com/terraform/language/modules/sources

	version = "3.0.0"
	source = "terraform-aws-modules/cloudfront/aws"
    tags_variable = "custom_tags"

  }
`

func TestAccAutocloudModule(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAutocloudModule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"autocloud_module.s3_bucket", "name", "S3Bucket"),
					resource.TestCheckResourceAttr(
						"autocloud_module.s3_bucket", "version", "3.0.0"),
					resource.TestCheckResourceAttr(
						"autocloud_module.s3_bucket", "source", "terraform-aws-modules/cloudfront/aws"),
					resource.TestCheckResourceAttrSet(
						"autocloud_module.s3_bucket", "variables.%"),
					resource.TestCheckResourceAttrSet(
						"autocloud_module.s3_bucket", "blueprint_config"),
					resource.TestCheckResourceAttr(
						"autocloud_module.s3_bucket", "variables.is_ipv6_enabled", "S3Bucket.is_ipv6_enabled"),
					resource.TestCheckResourceAttr(
						"autocloud_module.s3_bucket", "tags_variable", "custom_tags"),
					resource.TestCheckResourceAttrSet(
						"autocloud_module.s3_bucket", "blueprint_config"),
				),
			},
		},
	})
}
