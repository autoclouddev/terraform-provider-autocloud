package autocloud_provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testDataSourceTerraformProcessor = `

resource "autocloud_module" "s3_bucket" {

	name = "S3Bucket"
	version = "3.4.0"
	source  = "terraform-aws-modules/s3-bucket/aws"

  }

  data "autocloud_terraform_processor" "s3_processor" {
	source_module_id = autocloud_module.s3_bucket.id
	omit_variables   = ["request_payer", "attach_deny_insecure_transport_policy", "putin_khuylo", "attach_policy", "control_object_ownership", "attach_lb_log_delivery_policy", "create_bucket", "restrict_public_buckets", "attach_elb_log_delivery_policy", "object_ownership", "attach_require_latest_tls_policy", "policy", "block_public_acls", "bucket", "acl", "block_public_policy", "object_lock_enabled", "force_destroy", "ignore_public_acls"]

	# bucket_prefix, acceleration_status, expected_bucket_owner => these vars are of 'shortText' type
	# attach_public_policy is of 'radio' type ('checkbox' types are similar to 'radio' types)

	# OVERRIDE VARIABLE EXAMPLES
	# - overriding bucket_prefix 'shortText' into 'radio'
	override_variable {
	  variable_name = "bucket_prefix"
	  display_name  = "bucket prefix (from override block)"
      helper_text   = "bucket prefix helper text (from override block)"
	  form_config {
		type = "radio"
		field_options {
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
	override_variable {
	  variable_name = "acceleration_status"
	  form_config {
		type = "checkbox"
		field_options {
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

	override_variable {
	  variable_name = "attach_public_policy"
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

// how to omit and override variables (as JSON)
var builderJson = compactJson(`{
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
  }`)

// the form with omitted and overridden variables (as JSON string)
var formJson = compactJson(`[
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

func TestTerraformProcessor(t *testing.T) {
	dataKey := "data.autocloud_terraform_processor.s3_processor"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTerraformProcessor,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						dataKey, "source_module_id"),
					resource.TestCheckResourceAttr(
						dataKey, "builder", builderJson),
					resource.TestCheckResourceAttr(
						dataKey, "form_config", formJson),
				),
			},
		},
	})
}

func TestTerraformProcessorAtLeastOneConfigBlocksValidationError(t *testing.T) {
	expectedError := `At least 1 "form_config" blocks are required`
	terraform := `data "autocloud_terraform_processor" "s3_processor" {
		override_variable {
		}
	  }`
	validateErrors(t, expectedError, terraform)
}

func TestTerraformProcessorTooManyFormConfigBlocksValidationError(t *testing.T) {
	expectedError := "Too many form_config blocks"
	terraform := `data "autocloud_terraform_processor" "s3_processor" {
		override_variable {
			form_config {
			}
			form_config {
			}
		}
	  }`
	validateErrors(t, expectedError, terraform)
}

func TestTerraformProcessorFieldOptionsIsRequiredForRadiosError(t *testing.T) {
	expectedError := "One field_options block is required"
	terraform := `data "autocloud_terraform_processor" "s3_processor" {
		source_module_id = "dummy"
		override_variable {
			variable_name = "bucket_prefix"
			form_config {
				type = "radio"
			}
		}
	  }`
	validateErrors(t, expectedError, terraform)
}

func TestTerraformProcessorFieldOptionsIsRequiredForCheckboxesError(t *testing.T) {
	expectedError := "One field_options block is required"
	terraform := `data "autocloud_terraform_processor" "s3_processor" {
		source_module_id = "dummy"
		override_variable {
			variable_name = "bucket_prefix"
			form_config {
				type = "checkbox"
			}
		}
	  }`
	validateErrors(t, expectedError, terraform)
}

func TestTerraformProcessorShortTextCanNotHaveOptionsError(t *testing.T) {
	expectedError := "ShortText variables can not have options"
	terraform := `data "autocloud_terraform_processor" "s3_processor" {
		source_module_id = "dummy"
		override_variable {
			variable_name = "bucket_prefix"
			form_config {
				type = "shortText"
				field_options{

				}
			}
		}
	  }`
	validateErrors(t, expectedError, terraform)
}

func TestTerraformProcessorIsRequiredValidationsCanNotHaveAValueError(t *testing.T) {
	expectedError := "'isRequired' validation rule can not have a value"
	terraform := `data "autocloud_terraform_processor" "s3_processor" {
		source_module_id = "dummy"
		override_variable {
			variable_name = "bucket_prefix"
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
