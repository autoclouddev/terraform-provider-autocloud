provider "autocloud" {
    ###
    # AutoCloud API Endpoint URL
    #
    # Sets the endpoint that Terraform will talk to in order to determine state.
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
    # or via explicit configuration:
    # token = ""
}
