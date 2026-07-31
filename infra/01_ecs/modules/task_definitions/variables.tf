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
}

variable "ecs_execution_role_arn" {
  type        = string
  description = "ARN of the role for the ECS task execution role"
}

variable "ecr_commitflow_repository_url" {
  type        = string
  description = "ECR repository URL for commitflow image"
}

variable "ssm_parameter_github_pat" {
  type        = string
  description = "ARN of the SSM parameter for GitHub PAT"
}

variable "ssm_parameter_gitlab_pat" {
  type        = string
  description = "ARN of the SSM parameter for GitLab PAT"
}
