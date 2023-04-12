# begin - paste to the top of the generated terraform code (config.tf.tpl)
terraform {
    required_version = "~> 1.1.0"
    required_providers {
    aws = {
        source  = "hashicorp/aws"
        version = "~> 4.0"
    }
    }
}

variable "account_num" {
    type        = string
    description = "Target AWS account number, mandatory"
}

variable "aws_region" {
    description = "AWS region"
    type        = string
}

variable "aws_role" {
    description = "AWS role to assume"
    type        = string
}

provider "aws" {
    region = var.aws_region
    # The following code is for using cross account assume role
    assume_role {
    role_arn = "arn:aws:iam::${var.account_num}:role/${var.aws_role}"
    }
}

# end - paste (config.tf.tpl)