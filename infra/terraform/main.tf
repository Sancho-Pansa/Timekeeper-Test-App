data "aws_vpc" "default" {
  default = true
}

// SSM - Systems Manager
data "aws_ssm_parameter" "amazon_linux" {
  name = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-x86_64"
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }

  filter {
    name   = "default-for-az"
    values = ["true"]
  }
}

resource "aws_security_group" "timekeeper" {
  name_prefix = "${var.project_name}-"
  description = "Security group for Timekeeper API"

  vpc_id = data.aws_vpc.default.id

  tags = {
    Name = "${var.project_name}-sg"
  }
}

resource "aws_vpc_security_group_ingress_rule" "api" {
  security_group_id = aws_security_group.timekeeper.id
  description       = "Timekeeper HTTP access"
  from_port         = 8080
  to_port           = 8080
  ip_protocol       = "tcp"
  cidr_ipv4         = var.allowed_cidr
}

resource "aws_vpc_security_group_egress_rule" "outbound" {
  security_group_id = aws_security_group.timekeeper.id
  description       = "Outbound IP address"

  ip_protocol = "-1"
  cidr_ipv4   = "0.0.0.0/0"
}

resource "aws_iam_role" "ec2" {
  name = "${var.project_name}-ec2-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"

    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecr" {
  role = aws_iam_role.ec2.name

  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_iam_role_policy_attachment" "ssm" {
  role       = aws_iam_role.ec2.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_instance_profile" "ec2" {
  name = "${var.project_name}-profile"
  role = aws_iam_role.ec2.name
}

resource "aws_instance" "timekeeper" {
  ami           = data.aws_ssm_parameter.amazon_linux.value
  instance_type = var.instance_type

  subnet_id                   = sort(data.aws_subnets.default.ids)[0]
  associate_public_ip_address = true

  vpc_security_group_ids = [
    aws_security_group.timekeeper.id
  ]

  iam_instance_profile = aws_iam_instance_profile.ec2.name

  user_data = file("${path.module}/user-data.sh")

  metadata_options {
    http_endpoint = "enabled"
    http_tokens   = "required"
  }
  root_block_device {
    volume_type = "gp3"
    volume_size = 12
    encrypted   = true
  }
}