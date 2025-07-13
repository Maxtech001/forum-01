# AWS Provider Configuration
variable "aws_region" {
  type    = string
  default = "us-east-1"
}

# EC2 Configuration
variable "ami_id" {
  type    = string
  default = "ami-053b0d53c279acc90" # Ubuntu 22.04 LTS (us-east-1)
}

variable "instance_type" {
  type    = string
  default = "t3.medium"
}
