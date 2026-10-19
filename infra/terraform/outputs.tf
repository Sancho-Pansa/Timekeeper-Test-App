output "instance_id" {
  description = "Created EC2 instance ID"
  value       = aws_instance.timekeeper.id
}

output "public_ip" {
  value = aws_instance.timekeeper.public_ip
}

output "security_group_id" {
  value = aws_security_group.timekeeper.id
}

output "instance_role_arn" {
  description = "IAM Role attached to EC2"
  value       = aws_iam_role.ec2.arn
}