package blueprint_config_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	acctest "gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/acctest"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
)

func findBlueprintConfigById(configs []blueprint_config.BluePrintConfig, id string) blueprint_config.BluePrintConfig {
	for _, v := range configs {
		if v.Id == id {
			return v
		}
	}

	return blueprint_config.BluePrintConfig{}
}

func TestAccBlueprintConfig_sourceValidation(t *testing.T) {
	var blueprintConfig blueprint_config.BluePrintConfig
	resourceName := "data.autocloud_blueprint_config.test"
	experimental := true
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acctest.CreateMuxFactories(experimental),
		Steps: []resource.TestStep{
			{
				Config: testAccBlueprintConfig_basicSource(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlueprintConfigExist(resourceName, &blueprintConfig),
					resource.TestCheckResourceAttrSet(
						resourceName, "source.%"),
					resource.TestCheckResourceAttrSet(
						resourceName, "blueprint_config"),
					resource.TestCheckResourceAttrSet(
						resourceName, "config"),
				),
			},
		},
	})
}

func testAccBlueprintConfig_basicSource() string {
	return `
resource "autocloud_module" "s3" {
	name    = "s3"
	source  = "terraform-aws-modules/s3-bucket/aws"
	version = "3.6.0"
}

resource "autocloud_module" "kms" {
	name    = "kms"
	source  = "terraform-aws-modules/kms/aws"
	version = "1.3.0"
}
data "autocloud_blueprint_config" "test" {
	source = {
		kms = autocloud_module.kms.blueprint_config
		s3  = autocloud_module.s3.blueprint_config
	  }
}
`
}

func TestAccBlueprintConfig_empty(t *testing.T) {
	resourceName := "data.autocloud_blueprint_config.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBlueprintConfig_empty(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						resourceName, "id"),
					resource.TestCheckResourceAttrSet(
						resourceName, "blueprint_config"),
				),
			},
		},
	})
}

func testAccBlueprintConfig_empty() string {
	return `
	data "autocloud_blueprint_config" "test" {}
`
}

func TestAccBlueprintConfig_createConfig(t *testing.T) {
	var formVariables []generator.FormShape
	resourceName := "autocloud_blueprint.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBlueprintConfig_createConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCorrectVariablesLength(resourceName, &formVariables),
					resource.TestCheckResourceAttrSet(
						resourceName, "id"),
					resource.TestCheckResourceAttrSet(
						resourceName, "config"),
				),
			},
		},
	})
}

func testAccBlueprintConfig_createConfig() string {
	return `
	resource "autocloud_module" "kms" {
		name    = "kms"
		source  = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/security/kms?ref=0.1.0"
	}
	data "autocloud_blueprint_config" "level1_1" {
		source = {
			kms = autocloud_module.kms.blueprint_config
		}
	}
	data "autocloud_blueprint_config" "level1_2" {
		source = {
			kms = autocloud_module.kms.blueprint_config
		}
	}
	data "autocloud_blueprint_config" "level2" {
		source = {
			level1_1 = data.autocloud_blueprint_config.level1_1.blueprint_config
			level1_2 = data.autocloud_blueprint_config.level1_1.blueprint_config
		}
	}
	resource "autocloud_blueprint" "test" {
		name = "complexTree"
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


		file {
			action = "CREATE"

			destination = "eks-cluster-{{clusterName}}.tf"
			variables = {
			clusterName = "EKSGenerator.clusterName"
			}
			modules = ["EKSGenerator"]
		}


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
		#config = data.autocloud_blueprint_config.level2.config
		config = data.autocloud_blueprint_config.level1_1.config
	}
`
}

func TestAccBlueprintConfig_OverrideVars(t *testing.T) {
	var formVariables []generator.FormShape
	omitted := []string{
		"request_payer",
		"attach_deny_insecure_transport_policy",
		"putin_khuylo",
		"attach_policy",
		"control_object_ownership",
		"attach_lb_log_delivery_policy",
		"create_bucket",
		"restrict_public_buckets",
		"attach_elb_log_delivery_policy",
		"object_ownership",
		"attach_require_latest_tls_policy",
		"policy",
		"block_public_acls",
		"bucket",
		"acl",
		"block_public_policy",
		"object_lock_enabled",
		"force_destroy",
		"ignore_public_acls"}

	overideVars := []string{
		`variable {
			name         = "bucket_prefix"
			display_name = "bucket prefix (from override block)"
			helper_text  = "bucket prefix helper text (from override block)"

			type = "radio"
			options {
				option {
					label   = "dev"
					value   = "some-dev-prefix"
					checked = false
				}
				option {
					label   = "nonprod"
					value   = "some-nonprod-prefix"
					checked = true
				}
				option {
					label = "prod"
					value = "some-prod-prefix"
				}
			}
			  validation_rule {
				rule          = "isRequired"
				error_message = "invalid"
			  }

		  }`,
	}

	resourceName := "data.autocloud_blueprint_config.s3_processor"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				//nolint
				Config: testAccBlueprintConfig_OverrideVariables(omitted, overideVars),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCorrectVariablesLength(resourceName, &formVariables),
					testAccCheckOmitCorrectness(omitted, &formVariables),
					resource.TestCheckResourceAttrSet(
						resourceName, "id"),
					resource.TestCheckResourceAttrSet(
						resourceName, "config"),
					testAccCheckOverrides(resourceName, overideVars),
				),
			},
		},
	})
}

func testAccBlueprintConfig_OverrideVariables(omitted []string, overrides []string) string {
	template := `

	resource "autocloud_module" "s3_bucket" {

		name = "S3Bucket"
		version = "3.4.0"
		source  = "terraform-aws-modules/s3-bucket/aws"

	}

	data "autocloud_blueprint_config" "s3_processor" {
		source = {
		  s3 = autocloud_module.s3_bucket.blueprint_config
		}
		omit_variables = [
		  %s
		]

		# bucket_prefix, acceleration_status, expected_bucket_owner => these vars are of 'shortText' type
		# attach_public_policy is of 'radio' type ('checkbox' types are similar to 'radio' types)

		# OVERRIDE VARIABLE EXAMPLES
		# - overriding bucket_prefix 'shortText' into 'radio'
		%s

		# - overriding acceleration_status 'shortText' into 'checkbox'
		variable {
		  name = "acceleration_status"

			type = "checkbox"
			options {
				option {
				label = "Option 1"
				value = "acceleration_status_1"

				}
				option {
				label   = "Option 2"
				value   = "acceleration_status_2"
				checked = true
				}
				option {
				label   = "Option 3"
				value   = "acceleration_status_3"
				checked = true
				}

		  }
		}

		# - NOT overriding expected_bucket_owner 'shortText' (it should be displayed as shortText)
		# ...

		# - overriding attach_public_policy 'radio' into 'shortText'
		/*
		variable {
		  name = "attach_public_policy"

			type = "shortText"
			validation_rule {
			  rule          = "regex"
			  value         = "^(yes|no)$"
			  error_message = "invalid. you should choose between 'yes' or 'no'"
			}

		}*/
	  }

	`
	omittedInStrings := ""
	overrideVars := ""
	for _, v := range omitted {
		omittedInStrings += fmt.Sprintf(`"%s",`, v)
	}
	for _, v := range overrides {
		overrideVars += v
	}

	return fmt.Sprintf(template, omittedInStrings, overrideVars)
}

func TestGenericBlueprintConfig(t *testing.T) {
	dataKey := "data.autocloud_blueprint_config.generic"

	testDataSourceBluenprintConfig := `



	data "autocloud_blueprint_config" "generic" {
		variable {
			name = "project_name"

			type = "shortText"
			validation_rule {
				rule          = "isRequired"
				error_message = "invalid name"
			}

		}

		variable {
			name = "env"
			display_name = "environment target"
			helper_text  = "environment target description"

			type = "radio"
			options {
				option {
					label = "dev"
					value = "dev"
					checked = true
				}
				option {
					label   = "nonprod"
					value   = "nonprod"
				}
				option {
					label   = "prod"
					value   = "prod"
				}
			}
			validation_rule {
				rule          = "isRequired"
				error_message = "invalid"
			}

		}
	}`

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBluenprintConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						dataKey, "config"),
				),
			},
		},
	})
}

func TestBlueprintConfigWhenValueIsSetAFormConfigCanNotBeSet(t *testing.T) {
	expectedError := blueprint_config.ErrSetValueInForm
	terraform := `data "autocloud_blueprint_config" "s3_processor" {
		variable {
			name = "bucket_prefix"
			value = "some dummy value"

			type = "shortText"
			validation_rule {
				rule          = "isRequired"
				error_message = "invalid"
				value		  = "some value"
			}

		}
	  }`
	acctest.ValidateErrors(t, expectedError, terraform)
}

func TestBlueprintConfigTooManyFormConfigBlocksValidationError(t *testing.T) {
	expectedError := blueprint_config.ErrOneBlockOptionsRequied
	terraform := `data "autocloud_blueprint_config" "s3_processor" {
		variable {
			options {
			}
			options {
			}
		}
	  }`
	acctest.ValidateErrors(t, expectedError, terraform)
}

func TestBlueprintConfigFieldOptionsIsRequiredForRadiosError(t *testing.T) {
	expectedError := blueprint_config.ErrOneBlockOptionsRequied
	terraform := `data "autocloud_blueprint_config" "s3_processor" {
		variable {
			name = "bucket_prefix"
			type = "radio"

		}
	  }`
	acctest.ValidateErrors(t, expectedError, terraform)
}

func TestBlueprintConfigFieldOptionsIsRequiredForCheckboxesError(t *testing.T) {
	expectedError := blueprint_config.ErrOneBlockOptionsRequied
	terraform := `data "autocloud_blueprint_config" "s3_processor" {
		variable {
			name = "bucket_prefix"
			type = "checkbox"

		}
	  }`
	acctest.ValidateErrors(t, expectedError, terraform)
}

func TestBlueprintConfigShortTextCanNotHaveOptionsError(t *testing.T) {
	expectedError := blueprint_config.ErrShortTextCantHaveOptions
	terraform := `data "autocloud_blueprint_config" "s3_processor" {

		variable {
			name = "bucket_prefix"
			type = "shortText"
			options{

			}
		}
	  }`
	acctest.ValidateErrors(t, expectedError, terraform)
}

func TestBlueprintConfigIsRequiredValidationsCanNotHaveAValueError(t *testing.T) {
	expectedError := blueprint_config.ErrIsRequiredCantHaveValue
	terraform := `data "autocloud_blueprint_config" "s3_processor" {
		variable {
			name = "bucket_prefix"

			type = "shortText"
			validation_rule {
				rule          = "isRequired"
				error_message = "invalid"
				value		  = "some value"
			}

		}
	  }`
	acctest.ValidateErrors(t, expectedError, terraform)
}

func TestBuildVariableContent(t *testing.T) {
	variableContentVars := map[string]blueprint_config.VariableContent{
		"radio": {
			DisplayName: "this is display",
			HelperText:  "helper text",
			FormConfig: blueprint_config.FormConfig{
				Type: "radio",
				ValidationRules: []blueprint_config.ValidationRule{
					{
						Rule:         "isRequired",
						ErrorMessage: "invalid",
					},
				},
				FieldOptions: []blueprint_config.FieldOption{
					{
						Label:   "dev",
						Value:   "some-dev-prefix",
						Checked: true,
					},
					{
						Label:   "prod",
						Value:   "some-prod-prefix",
						Checked: false,
					},
					{
						Label:   "non-prod",
						Value:   "some-non-prod-prefix",
						Checked: false,
					},
				},
			},
		},
		"inputText": {
			DisplayName: "this is display",
			HelperText:  "helper text",
			FormConfig: blueprint_config.FormConfig{
				Type: "shortText",
				ValidationRules: []blueprint_config.ValidationRule{
					{
						Rule:         "isRequired",
						ErrorMessage: "invalid",
					},
				},
				FieldOptions: []blueprint_config.FieldOption{},
			},
		},
		"raw": {
			DisplayName: "this is a raw variable",
			HelperText:  "helper text",
			Value:       "var.ami",
			FormConfig: blueprint_config.FormConfig{
				Type: "raw",
			},
		},
		// "map": {
		// 	DisplayName: "this is display",
		// 	HelperText:  "helper text",
		// 	FormConfig: blueprint_config.FormConfig{
		// 		Type: "map",
		// 		ValidationRules: []blueprint_config.ValidationRule{
		// 			{
		// 				Rule:         "isRequired",
		// 				ErrorMessage: "invalid",
		// 			},
		// 		},
		// 		FieldOptions: []blueprint_config.FieldOption{},
		// 	},
		// 	//SHOULD THAT BE ARRAY??
		// 	RequiredValues: `[{
		// 		"arn": "dummy-ecs-cluster-arn-value",
		// 		"name": "dummy-ecs-cluster-arn-name"
		// 	}]`,
		// },
	}
	//this part of the code is bulding the test tables
	type test struct {
		input map[string]interface{}
		want  blueprint_config.VariableContent
	}
	tests := make(map[string]test)
	for testName, vc := range variableContentVars {
		tests[testName] = test{
			input: map[string]interface{}{
				"variable": []interface{}{
					// map[string]interface{}{
					// 	"name":         "great_name", // out of the scope of this test
					// },
					createRawVariableContentSchema(vc),
				},
			},
			want: vc,
		}
	}
	// Helps in sorting the labels for field options, so the order is correct in cmp.Diff
	trans := cmp.Transformer("SortFieldOptionsByLabel", func(in blueprint_config.VariableContent) blueprint_config.VariableContent {
		sort.Slice(in.FormConfig.FieldOptions, func(i, j int) bool {
			return in.FormConfig.FieldOptions[i].Label < in.FormConfig.FieldOptions[j].Label
		})
		return in
	})
	//running the test tables
	for name, tc := range tests {
		blueprintConfigDataSource := blueprint_config.DataSourceBlueprintConfig()
		testDataBlueprintResourceSchema := blueprintConfigDataSource.Schema
		d := schema.TestResourceDataRaw(t, testDataBlueprintResourceSchema, tc.input)
		v, _ := d.GetOk("variable")
		varsList := v.([]interface{})

		for _, currentVar := range varsList {
			varOverrideMap := currentVar.(map[string]interface{})
			got, _ := blueprint_config.BuildVariableFromSchema(varOverrideMap, &blueprint_config.BluePrintConfig{})

			if diff := cmp.Diff(&tc.want, got, trans); diff != "" {
				t.Fatalf("TESTNAME: %s BuildVariableFromSchema() mismatch (-want +got):\n%s", name, diff)
			}
		}
	}
}

// creates a raw schema from a VariableContent, mostly used for testing
func createRawVariableContentSchema(content blueprint_config.VariableContent) map[string]interface{} {
	ud := make(map[string]interface{})
	ud["display_name"] = content.DisplayName
	ud["helper_text"] = content.HelperText
	ud["type"] = content.FormConfig.Type
	ud["value"] = content.Value
	ud["default"] = content.Default
	ud["variables"] = content.Variables
	vRules := make([]interface{}, 0)
	for _, rule := range content.FormConfig.ValidationRules {
		rawRule := map[string]interface{}{
			"rule":          rule.Rule,
			"error_message": rule.ErrorMessage,
			"value":         rule.Value,
		}
		vRules = append(vRules, rawRule)
	}
	ud["validation_rule"] = vRules
	options := make([]interface{}, 0)
	for _, option := range content.FormConfig.FieldOptions {
		op := map[string]interface{}{
			"label":   option.Label,
			"value":   option.Value,
			"checked": option.Checked,
		}
		options = append(options, op)
	}
	ud["options"] = []interface{}{
		map[string]interface{}{
			"option": options,
		},
	}
	ud["required_values"] = content.RequiredValues

	return ud
}
