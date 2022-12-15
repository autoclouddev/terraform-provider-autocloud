package autocloud_provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestBlueprintConfigOverrideVariables(t *testing.T) {
	dataKey := "data.autocloud_blueprint_config.s3_processor"

	testDataSourceBluenprintConfig := `

	resource "autocloud_module" "s3_bucket" {

		name = "S3Bucket"
		version = "3.4.0"
		source  = "terraform-aws-modules/s3-bucket/aws"

	}

	data "autocloud_blueprint_config" "s3_processor" {
		source = {
			s3 = autocloud_module.s3_bucket.blueprint_config
		}
		omit_variables   = ["request_payer", "attach_deny_insecure_transport_policy", "putin_khuylo", "attach_policy", "control_object_ownership", "attach_lb_log_delivery_policy", "create_bucket", "restrict_public_buckets", "attach_elb_log_delivery_policy", "object_ownership", "attach_require_latest_tls_policy", "policy", "block_public_acls", "bucket", "acl", "block_public_policy", "object_lock_enabled", "force_destroy", "ignore_public_acls"]

		# bucket_prefix, acceleration_status, expected_bucket_owner => these vars are of 'shortText' type
		# attach_public_policy is of 'radio' type ('checkbox' types are similar to 'radio' types)

		# OVERRIDE VARIABLE EXAMPLES
		# - overriding bucket_prefix 'shortText' into 'radio'
		variable {
		name = "bucket_prefix"
		display_name  = "bucket prefix (from override block)"
		helper_text   = "bucket prefix helper text (from override block)"
		form_config {
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
		}
		}

		# - overriding acceleration_status 'shortText' into 'checkbox'
		variable {
		name = "acceleration_status"
		form_config {
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
		}

		# - NOT overriding expected_bucket_owner 'shortText' (it should be displayed as shortText)
		# ...

		# - overriding attach_public_policy 'radio' into 'shortText'

		variable {
		name = "attach_public_policy"
		form_config {
			type = "shortText"
			validation_rule {
			rule          = "regex"
			value         = "^(yes|no)$"
			error_message = "invalid. you should choose between 'yes' or 'no'"
			}
		}
		}
	}

	`
	/*
			// how to omit and override variables (as JSON)
			var builderJsonOverrideVariables = compactJson(`{
			"omitVariables": [
			  "attach_require_latest_tls_policy",
			  "ignore_public_acls",
			  "attach_policy",
			  "bucket",
			  "attach_deny_insecure_transport_policy",
			  "block_public_acls",
			  "attach_elb_log_delivery_policy",
			  "restrict_public_buckets",
			  "create_bucket",
			  "acl",
			  "policy",
			  "request_payer",
			  "attach_lb_log_delivery_policy",
			  "object_lock_enabled",
			  "force_destroy",
			  "object_ownership",
			  "control_object_ownership",
			  "block_public_policy",
			  "putin_khuylo"
			],
			"overrideVariable": {
			  "acceleration_status": {
				"variableName": "acceleration_status",
				"value": null,
				"displayName": "",
		        "helperText": "",
				"formConfig": {
				  "type": "checkbox",
				  "fieldOptions": [
					{
					  "label": "Option 1",
					  "value": "acceleration_status_1",
					  "checked": false
					},
					{
					  "label": "Option 2",
					  "value": "acceleration_status_2",
					  "checked": true
					},
					{
					  "label": "Option 3",
					  "value": "acceleration_status_3",
					  "checked": true
					}
				  ],
				  "validationRules": []
				}
			  },
			  "attach_public_policy": {
				"variableName": "attach_public_policy",
				"value": null,
				"displayName": "",
		        "helperText": "",
				"formConfig": {
				  "type": "shortText",
				  "fieldOptions": null,
				  "validationRules": [
					{
					  "rule": "regex",
					  "value": "^(yes|no)$",
					  "errorMessage": "invalid. you should choose between 'yes' or 'no'"
					}
				  ]
				}
			  },
			  "bucket_prefix": {
				"variableName": "bucket_prefix",
				"value": null,
				"displayName": "bucket prefix (from override block)",
		        "helperText": "bucket prefix helper text (from override block)",
				"formConfig": {
				  "type": "radio",
				  "fieldOptions": [
					{
					  "label": "dev",
					  "value": "some-dev-prefix",
					  "checked": false
					},
					{
					  "label": "nonprod",
					  "value": "some-nonprod-prefix",
					  "checked": true
					},
					{
					  "label": "prod",
					  "value": "some-prod-prefix",
					  "checked": false
					}
				  ],
				  "validationRules": [
					{
					  "rule": "isRequired",
					  "value": "",
					  "errorMessage": "invalid"
					}
				  ]
				}
			  }
			}
		  }`)*/
	/*
			// the form with omitted and overridden variables (as JSON string)
			var formJsonOverrideVariables = compactJson(`[
			{
			  "id": "S3Bucket.acceleration_status",
			  "type": "checkbox",
			  "module": "S3Bucket",
			  "formQuestion": {
				"fieldId": "S3Bucket.acceleration_status",
				"fieldType": "checkbox",
				"fieldLabel": "acceleration_status",
				"explainingText": "(Optional) Sets the accelerate configuration of an existing bucket. Can be Enabled or Suspended.",
				"fieldOptions": [
				  {
					"label": "Option 1",
					"fieldId": "S3Bucket.acceleration_status-acceleration_status_1",
					"value": "acceleration_status_1",
					"checked": false
				  },
				  {
					"label": "Option 2",
					"fieldId": "S3Bucket.acceleration_status-acceleration_status_2",
					"value": "acceleration_status_2",
					"checked": true
				  },
				  {
					"label": "Option 3",
					"fieldId": "S3Bucket.acceleration_status-acceleration_status_3",
					"value": "acceleration_status_3",
					"checked": true
				  }
				],
				"validationRules": []
			  },
			  "fieldDataType": "string",
			  "fieldDefaultValue": "null",
			  "fieldValue": "null"
			},
			{
			  "id": "S3Bucket.attach_public_policy",
			  "type": "shortText",
			  "module": "S3Bucket",
			  "formQuestion": {
				"fieldId": "S3Bucket.attach_public_policy",
				"fieldType": "shortText",
				"fieldLabel": "attach_public_policy",
				"explainingText": "Controls if a user defined public bucket policy will be attached (set to ` + "`false`" + ` to allow upstream to apply defaults to the bucket)",
				"fieldOptions": null,
				"validationRules": [
				  {
					"rule": "regex",
					"value": "^(yes|no)$",
					"errorMessage": "invalid. you should choose between 'yes' or 'no'"
				  }
				]
			  },
			  "fieldDataType": "bool",
			  "fieldDefaultValue": "true",
			  "fieldValue": "true"
			},
			{
			  "id": "S3Bucket.bucket_prefix",
			  "type": "radio",
			  "module": "S3Bucket",
			  "formQuestion": {
				"fieldId": "S3Bucket.bucket_prefix",
				"fieldType": "radio",
				"fieldLabel": "bucket prefix (from override block)",
				"explainingText": "bucket prefix helper text (from override block)",
				"fieldOptions": [
				  {
					"label": "dev",
					"fieldId": "S3Bucket.bucket_prefix-some-dev-prefix",
					"value": "some-dev-prefix",
					"checked": false
				  },
				  {
					"label": "nonprod",
					"fieldId": "S3Bucket.bucket_prefix-some-nonprod-prefix",
					"value": "some-nonprod-prefix",
					"checked": true
				  },
				  {
					"label": "prod",
					"fieldId": "S3Bucket.bucket_prefix-some-prod-prefix",
					"value": "some-prod-prefix",
					"checked": false
				  }
				],
				"validationRules": [
				  {
					"rule": "isRequired",
					"value": "",
					"errorMessage": "invalid"
				  }
				]
			  },
			  "fieldDataType": "string",
			  "fieldDefaultValue": "null",
			  "fieldValue": "null"
			},
			{
			  "id": "S3Bucket.expected_bucket_owner",
			  "type": "string",
			  "module": "S3Bucket",
			  "formQuestion": {
				"fieldId": "S3Bucket.expected_bucket_owner",
				"fieldType": "shortText",
				"fieldLabel": "expected_bucket_owner",
				"explainingText": "The account ID of the expected bucket owner",
				"fieldOptions": null,
				"validationRules": null
			  },
			  "fieldDataType": "string",
			  "fieldDefaultValue": "null",
			  "fieldValue": "null"
			}
		  ]`)
	*/
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBluenprintConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						dataKey, "source.%"),
					/*resource.TestCheckResourceAttr(
						dataKey, "builder", builderJsonOverrideVariables),
					resource.TestCheckResourceAttr(
						dataKey, "form_config", formJsonOverrideVariables),*/
				),
			},
		},
	})
}

func TestBlueprintConfigComposability(t *testing.T) {
	dataKey := "data.autocloud_blueprint_config.cf_processor"

	testDataSourceBluenprintConfigComposability := `

	resource "autocloud_module" "s3_bucket" {

		name = "S3Bucket"
		version = "3.4.0"
		source  = "terraform-aws-modules/s3-bucket/aws"

	}

	resource "autocloud_module" "cloudfront" {

		name = "Cloudfront"
		version = "3.0.0"
		source = "terraform-aws-modules/cloudfront/aws"

	}

	data "autocloud_blueprint_config" "s3_processor" {
		source = {
			s3 = autocloud_module.s3_bucket.blueprint_config
		}
		omit_variables = [
			"acceleration_status",
			"acl",
			"attach_deny_insecure_transport_policy",
			"attach_elb_log_delivery_policy",
			"attach_lb_log_delivery_policy",
			"attach_policy",
			"attach_public_policy",
			"attach_require_latest_tls_policy",
			"block_public_acls",
			"block_public_policy",
			"bucket_prefix",
			"control_object_ownership",
			"create_bucket",
			"expected_bucket_owner",
			"force_destroy",
			"ignore_public_acls",
			"object_lock_enabled",
			"object_ownership",
			"policy",
			"putin_khuylo",
			"request_payer",
		]
	}

	data "autocloud_blueprint_config" "cf_processor" {
		source = {
			cloudfront = autocloud_module.cloudfront.blueprint_config
		}

		# omitting most of the variables to simplify the form
		omit_variables = [
		"aliases",
		// "comment", // we'll take the value from s3 bucket name
		"create_distribution",
		"create_monitoring_subscription",
		"create_origin_access_identity",
		"custom_error_response",
		"default_cache_behavior",
		"default_root_object",
		"enabled",
		"geo_restriction",
		"http_version",
		"is_ipv6_enabled",
		"logging_config",
		"ordered_cache_behavior",
		"origin",
		"origin_access_identities",
		"origin_group",
		"price_class",
		"realtime_metrics_subscription_status",
		"retain_on_delete",
		"tags",
		"viewer_certificate",
		"wait_for_deployment",
		"web_acl_id",
		]

		# OVERRIDE VARIABLE EXAMPLES
		# - set values from other modules outputs
		variable {
		name = "comment"
		value         = autocloud_module.s3_bucket.outputs["s3_bucket_id"]
		}
	}

	`
	/*
			// how to omit and override variables (as JSON)
			var builderJsonComposability = compactJson(`{
			"omitVariables": [
			  "aliases",
			  "origin_access_identities",
			  "price_class",
			  "viewer_certificate",
			  "custom_error_response",
			  "http_version",
			  "create_origin_access_identity",
			  "create_distribution",
			  "enabled",
			  "origin",
			  "create_monitoring_subscription",
			  "retain_on_delete",
			  "tags",
			  "default_cache_behavior",
			  "geo_restriction",
			  "logging_config",
			  "wait_for_deployment",
			  "web_acl_id",
			  "is_ipv6_enabled",
			  "ordered_cache_behavior",
			  "realtime_metrics_subscription_status",
			  "default_root_object",
			  "origin_group"
			],
			"overrideVariable": {
			  "comment": {
				"variableName": "comment",
				"value": "module.S3Bucket.outputs.s3_bucket_id",
				"displayName": "",
				"helperText": "",
				"formConfig": {
				  "type": "",
				  "fieldOptions": null,
				  "validationRules": null
				}
			  }
			}
		  }`)

			// the form with omitted and overridden variables (as JSON string)
			var formJsonComposability = compactJson(`[
			{
			  "id": "Cloudfront.comment",
			  "type": "",
			  "module": "Cloudfront",
			  "formQuestion": {
				"fieldId": "Cloudfront.comment",
				"fieldType": "",
				"fieldLabel": "comment",
				"explainingText": "Any comments you want to include about the distribution.",
				"fieldOptions": null,
				"validationRules": []
			  },
			  "fieldDataType": "hcl-expression",
			  "fieldDefaultValue": "module.S3Bucket.outputs.s3_bucket_id",
			  "fieldValue": "module.S3Bucket.outputs.s3_bucket_id"
			}
		  ]`)
	*/
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBluenprintConfigComposability,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						dataKey, "source.%"),
					/*
						resource.TestCheckResourceAttr(
							dataKey, "builder", builderJsonComposability),
						resource.TestCheckResourceAttr(
							dataKey, "form_config", formJsonComposability),
					*/
				),
			},
		},
	})
}

func TestBlueprintConfigAtLeastOneConfigBlocksValidationError(t *testing.T) {
	expectedError := `A form_config must be defined for variable`
	terraform := `data "autocloud_blueprint_config" "s3_processor" {
		source_module_id = "dummy"
		variable {
			name = "bucket_prefix"
		}
	  }`
	validateErrors(t, expectedError, terraform)
}

func TestBlueprintConfigTooManyFormConfigBlocksValidationError(t *testing.T) {
	expectedError := "Too many form_config blocks"
	terraform := `data "autocloud_blueprint_config" "s3_processor" {
		variable {
			form_config {
			}
			form_config {
			}
		}
	  }`
	validateErrors(t, expectedError, terraform)
}

func TestBlueprintConfigFieldOptionsIsRequiredForRadiosError(t *testing.T) {
	expectedError := "One options block is required"
	terraform := `data "autocloud_blueprint_config" "s3_processor" {
		source_module_id = "dummy"
		variable {
			name = "bucket_prefix"
			form_config {
				type = "radio"
			}
		}
	  }`
	validateErrors(t, expectedError, terraform)
}

func TestBlueprintConfigFieldOptionsIsRequiredForCheckboxesError(t *testing.T) {
	expectedError := "One options block is required"
	terraform := `data "autocloud_blueprint_config" "s3_processor" {
		source_module_id = "dummy"
		variable {
			name = "bucket_prefix"
			form_config {
				type = "checkbox"
			}
		}
	  }`
	validateErrors(t, expectedError, terraform)
}

func TestBlueprintConfigShortTextCanNotHaveOptionsError(t *testing.T) {
	expectedError := "ShortText variables can not have options"
	terraform := `data "autocloud_blueprint_config" "s3_processor" {
		source_module_id = "dummy"
		variable {
			name = "bucket_prefix"
			form_config {
				type = "shortText"
				options{

				}
			}
		}
	  }`
	validateErrors(t, expectedError, terraform)
}

func TestBlueprintConfigIsRequiredValidationsCanNotHaveAValueError(t *testing.T) {
	expectedError := "'isRequired' validation rule can not have a value"
	terraform := `data "autocloud_blueprint_config" "s3_processor" {
		source_module_id = "dummy"
		variable {
			name = "bucket_prefix"
			form_config {
				type = "shortText"
				validation_rule {
					rule          = "isRequired"
					error_message = "invalid"
					value		  = "some value"
				  }
			}
		}
	  }`
	validateErrors(t, expectedError, terraform)
}

func TestBlueprintConfigWhenValueIsSetAFormConfigCanNotBeSet(t *testing.T) {
	expectedError := "A form_config can not be added when setting the variable's value."
	terraform := `data "autocloud_blueprint_config" "s3_processor" {
		source_module_id = "dummy"
		variable {
			name = "bucket_prefix"
			value = "some dummy value"
			form_config {
				type = "shortText"
				validation_rule {
					rule          = "isRequired"
					error_message = "invalid"
					value		  = "some value"
				  }
			}
		}
	  }`
	validateErrors(t, expectedError, terraform)
}

func TestGetFormBuilder(t *testing.T) {
	schema1 := dataSourceBlueprintConfig()
	testDataBlueprintResourceSchema := schema1.Schema
	fmt.Println(testDataBlueprintResourceSchema)
	raw := map[string]interface{}{
		"source": map[string]interface{}{
			"kms": `{
				"id": "clbnr5y2019144hyi1xc0yhex",
				"variables": [
					{
						"id": "kms.deletion_window_in_days",
						"type": "string",
						"module": "kms",
						"fieldValue": "",
						"formQuestion": {
							"fieldId": "kms.deletion_window_in_days",
							"fieldType": "shortText",
							"fieldLabel": "deletion_window_in_days",
							"fieldOptions": null,
							"explainingText": "KMS key deletion window in days",
							"validationRules": null
						},
						"fieldDataType": "number",
						"fieldDefaultValue": "10",
						"allowConsumerToEdit": true
					},
					{
						"id": "kms.description",
						"type": "string",
						"module": "kms",
						"fieldValue": "",
						"formQuestion": {
							"fieldId": "kms.description",
							"fieldType": "shortText",
							"fieldLabel": "description",
							"fieldOptions": null,
							"explainingText": "KMS key description indicating use case",
							"validationRules": null
						},
						"fieldDataType": "string",
						"fieldDefaultValue": "",
						"allowConsumerToEdit": true
					},
					{
						"id": "kms.enable_key_rotation",
						"type": "string",
						"module": "kms",
						"fieldValue": "",
						"formQuestion": {
							"fieldId": "kms.enable_key_rotation",
							"fieldType": "radio",
							"fieldLabel": "enable_key_rotation",
							"fieldOptions": [
								{
									"label": "Yes",
									"value": "true",
									"checked": true,
									"fieldId": "kms.enable_key_rotation-true"
								},
								{
									"label": "No",
									"value": "false",
									"checked": false,
									"fieldId": "kms.enable_key_rotation-false"
								}
							],
							"explainingText": "Whether or not AWS managed key rotation is enabled for this KMS key, defaults to true, enabled",
							"validationRules": null
						},
						"fieldDataType": "bool",
						"fieldDefaultValue": "true",
						"allowConsumerToEdit": true
					},
					{
						"id": "kms.enabled",
						"type": "string",
						"module": "kms",
						"fieldValue": "",
						"formQuestion": {
							"fieldId": "kms.enabled",
							"fieldType": "radio",
							"fieldLabel": "enabled",
							"fieldOptions": [
								{
									"label": "Yes",
									"value": "true",
									"checked": true,
									"fieldId": "kms.enabled-true"
								},
								{
									"label": "No",
									"value": "false",
									"checked": false,
									"fieldId": "kms.enabled-false"
								}
							],
							"explainingText": "Whether or not to create this resource, defaults to true, enabled",
							"validationRules": null
						},
						"fieldDataType": "bool",
						"fieldDefaultValue": "true",
						"allowConsumerToEdit": true
					},
					{
						"id": "kms.environment",
						"type": "string",
						"module": "kms",
						"fieldValue": "",
						"formQuestion": {
							"fieldId": "kms.environment",
							"fieldType": "shortText",
							"fieldLabel": "environment",
							"fieldOptions": null,
							"explainingText": "Environment KMS key belongs to",
							"validationRules": null
						},
						"fieldDataType": "string",
						"fieldDefaultValue": "null",
						"allowConsumerToEdit": true
					},
					{
						"id": "kms.name",
						"type": "string",
						"module": "kms",
						"fieldValue": "",
						"formQuestion": {
							"fieldId": "kms.name",
							"fieldType": "shortText",
							"fieldLabel": "name",
							"fieldOptions": null,
							"explainingText": "KMS key name",
							"validationRules": null
						},
						"fieldDataType": "string",
						"fieldDefaultValue": "null",
						"allowConsumerToEdit": true
					},
					{
						"id": "kms.namespace",
						"type": "string",
						"module": "kms",
						"fieldValue": "",
						"formQuestion": {
							"fieldId": "kms.namespace",
							"fieldType": "shortText",
							"fieldLabel": "namespace",
							"fieldOptions": null,
							"explainingText": "Namespace KMS key belongs to",
							"validationRules": null
						},
						"fieldDataType": "string",
						"fieldDefaultValue": "null",
						"allowConsumerToEdit": true
					}
				],
				"children": []
			  }`,
		},
	}

	d := schema.TestResourceDataRaw(t, testDataBlueprintResourceSchema, raw)
	formBuilder, err := getFormBuilder(d)
	if err != nil {
		fmt.Print(err)
		t.Fail()
	}
	variablesLength := len(formBuilder.BluePrintConfig.Variables)
	if variablesLength != 0 {
		t.Errorf("BlueprintConfig variables.length is not 0 is: %d", variablesLength)
	}
	nestedVariables := formBuilder.BluePrintConfig.Children[0].Variables
	if len(nestedVariables) != 7 {
		t.Errorf("BlueprintConfig.children[0].Variables.length is not 7 is: %d", len(nestedVariables))
	}
	fmt.Println(formBuilder.BluePrintConfig)
}
