data "autocloud_github_repos" "repos" {}

####
# Local variables
locals {
  # Destination repos where generated code will be submitted
  dest_repos = [
    for repo in data.autocloud_github_repos.repos.data[*].url : repo
    if length(regexall("/infrastructure-live-demo", repo)) > 0 || length(regexall("/self-hosted-infrastructure-live", repo)) > 0
  ]
}



####
# Module Resources
#
# Connect to the Terraform modules that will be used to create this generator


####
# KMS Key
#
resource "autocloud_module" "kms_key" {
  name   = "kmskey"
  source = "git@github.com:autoclouddev/infrastructure-modules-demo.git//aws/security/kms/key?ref=0.10.2"
}

data "autocloud_blueprint_config" "kms_key_processor" {
  source = {
    kms = autocloud_module.kms_key.blueprint_config
  }

  omit_variables = [
    # Use defaults in the module (don't collect)
    "customer_master_key_spec",
    "policies",

    # Hard Coded
    "enable_key_rotation",
  ]

  ###
  # Set description
  variable {
    name = "kms.variables.description"
    display_name = "Description"
    helper_text = "Description for key that appears in AWS console"

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

  ###
  # Set key deletion window
  variable {
    name         = "kms.variables.deletion_window_in_days"
    display_name = "Deletion Window"
    helper_text  = "Number of days to wait before deleting key permanently, defaults to 10 days"
    
    type = "shortText"
  }

  ###
  # Choose regionality
  variable {
    name = "kms.variables.multi_region"
    display_name = "Multi Region Key"
    helper_text  = "Whether or not the KMS key will be deployed as a multi region key"

    type = "radio"

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

  ###
  # Choose key usage
  variable {
    name = "kms.variables.key_usage"
    display_name = "Key Usage"

    type = "radio"

    options {
      option {
        label   = "Symmetric Encrypt/Decrypt"
        value   = "ENCRYPT_DECRYPT"
        checked = true
      }
      option {
        label   = "Signing/Verification"
        value   = "SIGN_VERIFY"
        checked = false
      }
      option {
        label   = "HMAC"
        value   = "GENERATE_VERIFY_MAC"
        checked = false
      }
    }
  }
}



####
# Create Blueprint Config
#
# Combine resources into the final config
data "autocloud_blueprint_config" "final" {
  source = {
    kms = data.autocloud_blueprint_config.kms_key_processor.blueprint_config
  }

  ###
  # Hide variables from user
  omit_variables = [
    # Global
    "enabled",
  ]

  ###
  # Hard code `enabled` to true to create all assets
  variable {
    name  = "enabled"
    value = true
  }

  ###
  # Set the namespace
  variable {
    name         = "namespace"
    display_name = "Namespace"
    helper_text  = "The organization namespace the assets will be deployed in"

    type = "shortText"

    value = "unstyl"
  }

  ###
  # Choose the environment
  variable {
    name         = "environment"
    display_name = "Environment"
    helper_text  = "The environment the assets will be deployed in"

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

    validation_rule {
      rule          = "isRequired"
      error_message = "invalid"
    }
  }

  ###
  # Input the name
  variable {
    name         = "name"
    display_name = "Name"
    helper_text  = "The name of the KMS key"

    type = "shortText"

    validation_rule {
      rule          = "isRequired"
      error_message = "This field is required"
    }
  }

  variable {
    name    = "tags"
    display_name = "Tags"
    helper_text  = "A map of tags to apply to the deployed assets"

    type = "map"

    # validation_rule {
    #   rule          = "isRequired"
    #   error_message = "invalid"
    # }
  }
}



####
# Create Blueprint
#
# Create generator blueprint that contains all the elements
resource "autocloud_blueprint" "this" {
  name = "KMS Key"

  ###
  # UI Configuration
  #
  author       = "chris@autocloud.dev"
  description  = "KMS Key"
  instructions = <<-EOT
  To deploy this generator, follow these simple steps:

  step 1: step-1-description
  step 2: step-2-description
  step 3: step-3-description
  EOT

  labels = ["aws"]

  ###
  # Form configuration
  config = data.autocloud_blueprint_config.final.config


  ###
  # Destination repository git configuraiton
  #
  # TODO:
  # - Reference site name using global values
  git_config {
    destination_branch = "main"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default

    pull_request {
      title                   = "[AutoCloud] new KMS Key {{namespace}}-{{environment}}-{{name}}, created by {{authorName}}"
      commit_message_template = "[AutoCloud] new KMS Key {{namespace}}-{{environment}}-{{name}}, created by {{authorName}}"
      body                    = file("./files/pull_request.md.tpl")
      variables = {
        authorName = "generic.authorName"
        namespace   = "kmskey.namespace"
        environment = "kmskey.environment"
        name        = "kmskey.name"
      }
    }
  }


  ###
  # File definitions
  #
  file {
    action      = "CREATE"
    destination = "aws/{{environment}}/{{namespace}}-{{environment}}-{{name}}.tf"
    variables = {
        namespace   = "kmskey.namespace"
        environment = "kmskey.environment"
        name        = "kmskey.name"
    }

    modules = [
      autocloud_module.kms_key.name
    ]
  }
}
