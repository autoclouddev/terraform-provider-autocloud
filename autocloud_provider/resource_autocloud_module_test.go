package autocloud_provider

import (
	"testing"

	"github.com/joho/godotenv"

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

  }
`

func TestAccAutocloudModule(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
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
						"autocloud_module.s3_bucket", "template"),
					resource.TestCheckResourceAttrSet(
						"autocloud_module.s3_bucket", "variables.%"),
					resource.TestCheckResourceAttr(
						"autocloud_module.s3_bucket", "variables.is_ipv6_enabled", "s3bucket.is_ipv6_enabled"),
				),
			},
		},
	})
}
