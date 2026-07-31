output "ssm_parameter_db_engine" {
  value = aws_ssm_parameter.db_engine.arn
}

output "ssm_parameter_db_address" {
  value = aws_ssm_parameter.db_address.arn
}

output "ssm_parameter_db_port" {
  value = aws_ssm_parameter.db_port.arn
}

output "ssm_parameter_db_name" {
  value = aws_ssm_parameter.db_name.arn
}

output "ssm_parameter_db_username" {
  value = aws_ssm_parameter.db_username.arn
}

output "db_secret_arn" {
  value = aws_db_instance.this.master_user_secret[0].secret_arn
}
