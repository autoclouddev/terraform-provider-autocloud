package blueprint_config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/acctest"
)

// this test verifies the correct values are read by .tf and provides examples of how to set the variable values
func TestBlueprintConfigDynamicTypeField(t *testing.T) {
	const dataKey = "data.autocloud_blueprint_config.test"

	testDataSourceBluenprintConfig := `
	  locals {
		localValue = "local-value"
	  }

	  data "autocloud_blueprint_config" "test" {
		// var #0
		variable {
		  name  = "dummy-string"
		  value = "this-is-a-string"
		}

		// var #1
		variable {
		  name  = "dummy-number"
		  value = 5
		}

		// var #2
		variable {
		  name  = "dummy-bool"
		  value = true
		}

		// var #3
		variable {
		  name = "dummy-list-string"
		  value = jsonencode([
			"my_policy_string"
		  ])
		}

		// var #4
		variable {
		  name = "dummy-list-number"
		  value = jsonencode([
			1234
		  ])
		}

		// var #5
		variable {
		  name = "dummy-list-string-with-var"
		  value = jsonencode([
			"my_policy_string",
			local.localValue
		  ])
		}

		// var #6
		variable {
		  name = "dummy-object"
		  value = jsonencode({
			my_tag : "my_tag_value"
		  })
		}
		// var #7 - tf syntax
		variable {
		  name = "dummy-nested-object"
		  value = jsonencode({
			managed-by = "autocloud" # Static value
			owner      = null        # Force user to enter
			"business-unit" = [      # Force user to choose
			  "finance",
			  "legal",
			  "engineering",
			  "sales"
			]
		  })
		}

	  }`

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBluenprintConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						dataKey, "variable.0.value", "this-is-a-string"),
					resource.TestCheckResourceAttr(
						dataKey, "variable.1.value", "5"),
					resource.TestCheckResourceAttr(
						dataKey, "variable.2.value", "true"),
					resource.TestCheckResourceAttr(
						dataKey, "variable.3.value", "[\"my_policy_string\"]"),
					resource.TestCheckResourceAttr(
						dataKey, "variable.4.value", "[1234]"),
					resource.TestCheckResourceAttr(
						dataKey, "variable.5.value", "[\"my_policy_string\",\"local-value\"]"),
					resource.TestCheckResourceAttr(
						dataKey, "variable.6.value", "{\"my_tag\":\"my_tag_value\"}"),
					resource.TestCheckResourceAttr(
						dataKey, "variable.7.value", "{\"business-unit\":[\"finance\",\"legal\",\"engineering\",\"sales\"],\"managed-by\":\"autocloud\",\"owner\":null}"),
				),
			},
		},
	})
}
