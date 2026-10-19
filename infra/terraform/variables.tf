variable "aws_region" {
  description = "Selected AWS region"
  type        = string
  default     = "eu-north-1"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"
}

variable "allowed_cidr" {
  description = "Allowed IPv4 CIDR for API"
  type        = string
  default     = "192.168.1.1/32"
}

variable "project_name" {
  description = "Resource name prefix"
  type        = string
  default     = "sancho-timekeeper-dev"
}