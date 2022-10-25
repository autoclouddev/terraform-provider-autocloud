package autocloud_provider

import (
	"testing"

	"github.com/joho/godotenv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccAutocloudModule1 = `
resource "autocloud_module_1" "s3_bucket" {

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

	version = "3.4.0"
	source = "terraform-aws-modules/s3-bucket/aws"

  }
`

func TestAccAutocloudModule1(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAutocloudModule1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"autocloud_module_1.s3_bucket", "name", "S3Bucket"),
					resource.TestCheckResourceAttr(
						"autocloud_module_1.s3_bucket", "version", "3.4.0"),
					resource.TestCheckResourceAttr(
						"autocloud_module_1.s3_bucket", "source", "terraform-aws-modules/s3-bucket/aws"),
					resource.TestCheckResourceAttrSet(
						"autocloud_module_1.s3_bucket", "template"),
					resource.TestCheckResourceAttrSet(
						"autocloud_module_1.s3_bucket", "variables"),
				),
			},
		},
	})
}
