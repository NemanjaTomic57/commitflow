##################################################
# VPC
##################################################

output "public_subnet_ids" {
  description = "IDs of the public subnets in the VPC"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "IDs of the private subnets in the VPC"
  value       = module.vpc.private_subnet_ids
}

##################################################
# Kafka cluster
##################################################

output "bastion_public_ip" {
  description = "Public IP address of the bastion"
  value       = module.kafka_cluster.bastion_public_ip
}

output "kafka_private_ips" {
  description = "Private IP addresses of Kafka nodes"
  value       = module.kafka_cluster.kafka_private_ips
}
