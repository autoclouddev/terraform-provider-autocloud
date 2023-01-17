# Provisional sales demo

This contains the tf file to provision a demo EKS generator

This has a bug, you have to manually paste the eks_demo.json file in the following DB models

InfrastructureChangeTemplate -> (select the deployed instance of main.tf) -> form

IacTerraformModule -> EKS instance -> Variables, please in here just paste `variables` of the whole json
