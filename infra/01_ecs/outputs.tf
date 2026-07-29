##################################################
# Terraform Remote State: Base
##################################################

output "kafka_private_ips" {
  description = "Private IP addresses of Kafka nodes"
  value       = data.terraform_remote_state.vpc.outputs.kafka_private_ips
}

##################################################
# ECS
##################################################

output "ecs_cluster_id" {
  description = "ID of the ECS cluster"
  value       = aws_ecs_cluster.commitflow.id
}
