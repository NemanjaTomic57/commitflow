variable "name" {
  type        = string
  description = "Name of the application"
}

variable "aws_region" {
  type = string
}

variable "ecs_task_role_arn" {
  type        = string
  description = "ARN of the role for the ECS task role"
  default     = "arn:aws:iam::761018874759:role/ECSCommitFlowTaskRole"
}

variable "ecs_execution_role_arn" {
  type        = string
  description = "ARN of the role for the ECS task execution role"
  default     = "arn:aws:iam::761018874759:role/ECSCommitFlowTaskExecutionRole"
}

variable "ssm_parameter_github_pat" {
  type        = string
  description = "ARN of the SSM parameter for GitHub PAT"
  default     = "arn:aws:ssm:eu-central-1:761018874759:parameter/commitflow/passwords/github-pat"
}

variable "ssm_parameter_gitlab_pat" {
  type        = string
  description = "ARN of the SSM parameter for GitLab PAT"
  default     = "arn:aws:ssm:eu-central-1:761018874759:parameter/commitflow/passwords/gitlab-pat"
}

variable "ecr_commitflow_repository_url" {
  type        = string
  description = "ECR repository URL for commitflow image"
}

variable "kafka_bootstrap_server" {
  type        = string
  description = "IP address for the Kafka broker endpoint"
}
