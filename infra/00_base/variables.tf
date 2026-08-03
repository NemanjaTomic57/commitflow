variable "db_password" {
  type      = string
  sensitive = true
}

variable "github_pat" {
  type      = string
  sensitive = true
}

variable "gitlab_pat" {
  type      = string
  sensitive = true
}
