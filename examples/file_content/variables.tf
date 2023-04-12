variable "account_num" {
  type        = string
  description = "Target AWS account number, mandatory"
}

variable "aws_role" {
  description = "AWS role to assume"
  type        = string
}
