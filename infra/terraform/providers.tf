provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "Sancho-Timekeeper"
      Environment = "dev"
      ManagedBy   = "Terraform"
    }
  }
}