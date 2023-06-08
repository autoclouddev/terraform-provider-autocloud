---
page_title: "autocloud_blueprint_config Data Source - AutoCloud"
subcategory: ""
description: |-
  terraform form processor (form builder)Creates, composes, and modifies Terraform blueprint form configurations, defining the form elements that will be shown to the end user.
---

# Data Source: autocloud_blueprint_config

Creates, composes, and modifies AutoCloud Terraform blueprint form configurations, defining the form elements that will be shown to the end user.

-> **Note**
The `autocloud_blueprint_config` data resource is the primary facility to define the user experience for the end user of the AutoCloud Terraform blueprint. See the [AutoCloud documentation](https://docs.autocloud.io/blueprint-configuration) for more details on how this resource generates forms to collect data from the user.

## Example Usage
Sample usage to deploy a KMS key to AWS:
```terraform
resource "autocloud_module" "kms_key" {
  name    = "cpkmskey"
  source  = "cloudposse/kms-key/aws"
  version = "0.12.1"
}

data "autocloud_blueprint_config" "kms_key_processor" {
  source = {
    kms = autocloud_module.kms_key.blueprint_config
  }

  omit_variables = [
    # Use Defaults
    "additional_tag_map",
    "alias",
    "attributes",
    "context",
    "customer_master_key_spec",
    "delimiter",
    "descriptor_formats",
    "id_length_limit",
    "key_usage",
    "label_key_case",
    "label_order",
    "label_value_case",
    "labels_as_tags",
    "regex_replace_chars",
    "stage",
    "tenant",

    # Hard Code
    "deletion_window_in_days",
    "enable_key_rotation",
    "enabled",
    "multi_region",
    "namespace",
  ]

  ###
  # Hard code `enabled` to true to create all assets
  variable {
    name  = "kms.variables.enabled"
    value = true
  }

  variable {
    name         = "kms.variables.namespace"
    display_name = "Namespace"
    helper_text  = "The organization namespace the new KMS key will be deployed in"

    type = "raw"

    value = "var.workspace_name"
  }

  ###
  # Choose the environment
  variable {
    name         = "kms.variables.environment"
    display_name = "Environment"
    helper_text  = "The environment the new KMS key will be deployed in"

    type = "radio"

    options {
      option {
        label   = "Nonprod"
        value   = "nonprod"
        checked = true
      }
      option {
        label = "Production"
        value = "production"
      }
    }
  }

  ###
  # Input the name
  variable {
    name         = "kms.variables.name"
    display_name = "Name"
    helper_text  = "The name of the KMS key"

    type = "shortText"

    validation_rule {
      rule          = "isRequired"
      error_message = "This field is required"
    }
  }

  ###
  # Set description
  variable {
    name = "kms.variables.description"
    display_name = "Key description"

    type = "shortText"
  }

  ###
  # Force key rotation
  variable {
    name  = "kms.variables.enable_key_rotation"
    display_name = "Automatic Key Rotation"
    helper_text  = "Enable automatic key rotation for the KMS key"

    value = true
  }

  variable {
    name = "kms.variables.deletion_window_in_days"
    type = "shortText"

    value = 14
  }

  variable {
    name         = "kms.variables.multi_region"
    display_name = "Multi Region Key"
    helper_text  = "Whether or not the KMS key will be deployed as a multi region key"
    type         = "radio"
    options {
      option {
        label   = "Single region key"
        value   = "true"
        checked = false
      }
      option {
        label   = "Multi region key"
        value   = "false"
        checked = true
      }
    }
  }

  variable {
    name         = "policy"
    display_name = "Key Policy"
    helper_text  = "The AWS KMS key policy to apply to the key"

    type = "editor"
  }

  variable {
    name    = "tags"
    display_name = "Tags"
    helper_text  = "A map of tags to apply to the deployed assets"

    type = "map"
  }
}

data "autocloud_blueprint_config" "final" {
  source = {
    kms = data.autocloud_blueprint_config.kms_key_processor.blueprint_config
  }

  omit_variables = [
    # Use Defaults

    # Hard Coded and Hidden from User
    "deletion_window_in_days",
    "enable_key_rotation",
    "enabled",
    "multi_region",
  ]

  display_order {
    priority = 0
    values = [
      "kms.variables.environment",
      "kms.variables.name",
      "kms.variables.description",
      "kms.variables.policy",
      "kms.variables.tags",
    ]
  }
}

resource "autocloud_blueprint" "this" {
  name = "CloudPosse AWS KMS Key"

  author       = "jim@example.com"
  description  = "Deploys a KMS Key using the CloudPosse KMS Key Module"
  instructions = <<-EOT
  To deploy this generator, these simple steps:

    * step 1: Choose the target environment
    * step 2: Provide a name to identify the KMS key
    * step 3: Provide a description for the key's intended use
    * step 4: Provide a key policy
    * step 3: Add tags to apply to key
  EOT

  labels = ["aws"]

  config = data.autocloud_blueprint_config.final.config


  file {
    action      = "CREATE"
    destination = "aws/{{environment}}/{{environment}}-{{name}}/main.tf"

    variables = {
      environment = "cpkmskey.environment"
      name        = "cpkmskey.name"
    }

    modules = [
      autocloud_module.kms_key.name
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

- [`display_order`](#nestedblock--display_order) - (Optional, Block List, Max: 1) Configuration block defining the order in which the variables will be displayed to the user. Defaults to alphanumeric order by the variable `name` value. Detailed below.
- `omit_variables` - (Optional, Set of String) A list of variables from the source blueprint configurations to omit from the blueprint form, preventing them from being shown to the user.
- `source` - (Optional, Map of String) A map of source blueprint configurations used as a starting point, which `display_order`, `omit_variables`, and `variable` blocks will modify.
- [`variable`](#nestedblock--variable) - (Optional, Block List) Configuration block defining or overriding a variable to create a form element to show to the user. Detailed below.

<a id="nestedblock--display_order"></a>
### `display_order`

The following arguments are supported:

- `priority` - (Number) The priority assigned to this display order. Used to resolve conflicts between display orders defined in blueprint configurations passed to the `autocloud_blueprint` resource. Higher priority indicates higher precedence.
- `values` - (List of String) A list of the variables to display, in the order they are to be displayed.

<a id="nestedblock--variable"></a>
### `variable`

The following arguments are supported:

- `name` - (String) The name of the variable to define or modify.
- [`conditional`](#nestedblock--variable--conditional) - (Optional, Block Set) Configuration block defining variable configurations that change based on the value of another variable's value. Detailed below.
- `display_name` - (Optional, String) A display name that will be shown to the user instead of the literal Terraform variable name.
- `helper_text` - (Optional, String) A string shown to the user providing additional detail and context to help them enter the correct information.
- [`options`](#nestedblock--variable--options) - (Optional, Block Set, Max: 1) Configuration block defining the options for a variable with a variable of type `radio`. Detailed below.
- `required_values` - (Optional, String) A json encoded string containing required values for `list` or `map` Terraform variable types.
- `type` - (Optional, String) The form display element used to collect the value of the variable from the user. Inferred from Terraform variable type where possible by default. Possible values are:
    - `checkbox` - Standard checkbox
    - `editor` - A code editor for entering long form text like boot scripts, policy documents, or source code.
    - `map` - A set of key-value pairs.
    - `radio` - Standard radio button.
    - `raw` - A raw string that will be written verbatim to the output.
    - `shortText` - Standard text input.

- `variables` (Optional, Map of String) A dictionary of keys and values to be used in AutoCloud string interpolation.
- [`validation_rule`](#nestedblock--variable--validation_rule) - (Optional, Block Set) Configuration block defining a validation rule to apply to the form element collecting the variable's value from the user. Detailed below.
- `value` - (Optional, String) The variable's value. References and AutoCloud string interpolation supported.

<a id="nestedblock--variable--conditional"></a>
### `conditional`

The following arguments are supported:

- `condition` - (String) The value of the source variable that will
- [`content`](#nestedblock--variable--conditional--content) - (Block Set, Min: 1, Max: 1) Configuration block defining the configuration of the variable when the condition is met. Detailed below.
- `source` - (String) A reference to the source variable whose value will determine the configuration of this variable.

<a id="nestedblock--variable--conditional--content"></a>
### `content`

The content of a `variable` block if a `conditional`'s condition is met.

-> **Note**
These arguments for the `content` block are identical to their counterparts in the `variable` block.

The following arguments are supported:

- `display_name` - (Optional, String) A display name that will be shown to the user instead of the literal Terraform variable name.
- `helper_text` - (Optional, String) A string shown to the user providing additional detail and context to help them enter the correct information.
- [`options`](#nestedblock--variable--options) - (Optional, Block Set, Max: 1) Configuration block defining the options for a variable with a variable of type `radio`. Detailed below.
- `required_values` - (Optional, String) A json encoded string containing required values for `list` or `map` Terraform variable types.
- `type` - (Optional, String) The form display element used to collect the value of the variable from the user. Inferred from Terraform variable type where possible by default. Possible values are:
    - `checkbox` - Standard checkbox
    - `editor` - A code editor for entering long form text like boot scripts, policy documents, or source code.
    - `map` - A set of key-value pairs.
    - `radio` - Standard radio button.
    - `raw` - A raw string that will be written verbatim to the output.
    - `shortText` - Standard text input.
- `variables` (Optional, Map of String) A dictionary of keys and values to be used in AutoCloud string interpolation.
- [`validation_rule`](#nestedblock--variable--validation_rule) - (Optional, Block Set) Configuration block defining a validation rule to apply to the form element collecting the variable's value from the user. Detailed below.
- `value` - (Optional, String) The variable's value. References and AutoCloud string interpolation supported.

<a id="nestedblock--variable--options"></a>
### `options`

The following arguments are supported:

- [`option`](#nestedblock--variable--options--option) - (Optional, Block Set) Configuration block defining a radio button form element option. Detailed below.

<a id="nestedblock--variable--options--option"></a>
### `option`

Defines a radio button form element option.

The following arguments are supported:

- `label` - (String) The display label shown to the user.
- `value` - (String) The value of the form element when this option is selected.
- `checked` - (Optional, Boolean) Whether or not the element is checked. Defaults to false, not checked.

<a id="nestedblock--variable--validation_rule"></a>
### `validation_rule`

Defines a validation rule to apply to the form element collecting the variable's value from the user. More than one validation rule may be defined for a variable.

-> **Note**
Not all validation rules apply to all variable types. For more information on AutoCloud blueprint form validation, see the [AutoCloud documentation](https://docs.autocloud.io/variable-definitions#-enBz_)

The following arguments are supported:

- `rule` - (String) The rule to apply to the field. Possible values are:
    - `isRequired` - Indicates a required field.
    - `maxLength` - The maximum number of elements in a list or map.
    - `minLength` - The minimum number of elements in a list or map.
    - `regex` - A regular expression to validate a text field, uses [RE2 syntax](https://github.com/google/re2/wiki/Syntax).
- `error_message` - (String) An error message displayed to the user when the rule is violated.
- `value` - (Optional, String) A value appropriate to the rule provided.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `blueprint_config` - (String) The autogenerated blueprint configuration that can be passed to an `autocloud_blueprint_config` data resource as a source.
- `config` - (String) The completed form configuration that an `autocloud_blueprint` resource will use to render the data collection form to the end user.
- `id` - (String) The unique identifier for the data resource.
