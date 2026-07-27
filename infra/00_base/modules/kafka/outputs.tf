output "bastion_public_ip" {
  value = aws_instance.bastion.public_ip
}

output "kafka_private_ips" {
  value = {
    for name, instance in aws_instance.kafka :
    name => instance.private_ip
  }
}

output "kafka_bootstrap_server" {
  value = join(",", [
    for _, instance in aws_instance.kafka :
    "${instance.private_ip}:9092"
  ])
}
