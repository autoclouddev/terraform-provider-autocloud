provider "autocloud" {
    ###
    # AutoCloud API Endpoint URL
    # 
    # Sets the endpoint that Terraform will talk to in order to determine state. Omit to use production SaaS endpoint.
    # For self hosted or testing environments, set either via environment variable:
    # 
    # export AUTOCLOUD_API=https://api.autocloud.io/api/v.0.0.1
    #
    # or via explicit configuration:

    # endpoint = "https://api.autocloud.io/api/v.0.0.1" # omit this, use autocloud prod for default



    ###
    # AutoCloud API Token
    #
    # Authorizes user to interact with AutoCloud API. These must be generated here:
    #
    # https://app.autocloud.io/settings/integrations/terraform-provider
    #
    # Value must be set eithe via environment variable:
    #
    # export AUTOCLOUD_TOKEN=
    # 
    # or via explicit configuraiton:

    # token = ""



    ###
    # AutoCloud Terraform Provider Version
    #
    # Set via normal Terraform version spec or omit for latest that matches version restrictions set in `required_providers`
    # configuration.

    # version  = "~> 0.4.0"
}
