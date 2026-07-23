variable "ec2_image_id" {
  description = "Debian Trixie (ARM64) golden AMI made with HashiCorp Packer"
  type        = string
  default     = "ami-0723bff07f72bb394"
}

variable "key_name" {
  description = "Existing EC2 Key Pair"
  type        = string
  default     = "aws"
}
