# terraform-provider

This project contains three main projects: 

+ `/sdk/autocloud_sdk` contains the go sdk, an api client for our current api

+ `/terraform_provider/autocloud_provider` contains the implementation of the Terraform Provider of autocloud

+ `terraform/autocloud` contains the producer flow example of the IAC catalog


You need to have go installed in your machine https://go.dev/doc/install
Also you will need terraform  https://learn.hashicorp.com/tutorials/terraform/install-cli


To setup the env variables run the following once your .env file is complete

`$ export $(grep -v '^#' .env | xargs)`

Also, you must log in to your sso aws account to grab the needed keys for cognito

make sure you have the profile in you local machine

```
# on file ~/.aws/credentials 
[autocloud-aws-sso-sandbox-developer]
sso_start_url = https://auto-cloud.awsapps.com/start
sso_region = us-east-1
sso_account_id = 632941798677
sso_role_name = universal-developer
```

Then run the following:

`$ aws sso login --profile autocloud-aws-sso-sandbox-developer`
