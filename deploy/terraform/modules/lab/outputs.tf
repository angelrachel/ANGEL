output "instance_id" {
  value       = aws_instance.lab.id
  description = "ID of the lab instance."
}

output "instance_ip" {
  value       = aws_instance.lab.public_ip
  description = "Public IP of the lab instance."
}
