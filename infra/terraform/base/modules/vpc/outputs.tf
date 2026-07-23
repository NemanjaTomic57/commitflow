output "bastion_public_ip" {
  description = "Public IP address of the bastion"
  value       = aws_instance.bastion.public_ip
}

output "kafka_private_ips" {
  description = "Private IP addresses of Kafka nodes"
  value = {
    for name, instance in aws_instance.kafka :
    name => instance.private_ip
  }
}
