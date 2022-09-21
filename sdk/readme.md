# autocloud sdk

This folder contains the sdk project which is communicating directly with the api
It needs the correct env variables and your aws sso session to work correctly


To run it go to cd ./sdk
Run the following commands
`go mod tidy`
`go mod vendor`

This setup has a multi module setup to test the sdk as if the terraform provider is calling it
All code is in main.go file

Return back to the sdk folder
`cd ../`

Setup env variables and aws sso login 

`$ export $(grep -v '^#' ../.env | xargs)`

`$ export AWS_PROFILE=autocloud-aws-sso-sandbox-developer`


Finally run go `run main.go` to test the sdk

Feel free to modify main.go to test the available commands
