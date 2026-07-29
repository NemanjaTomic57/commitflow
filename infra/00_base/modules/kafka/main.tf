resource "aws_instance" "kafka" {
  for_each = var.private_subnet_ids

  ami                    = var.ami_id
  instance_type          = var.instance_type
  key_name               = var.key_name
  subnet_id              = each.value
  vpc_security_group_ids = [var.kafka_security_group_id]

  tags = {
    Name = "${var.name}-kafka-${each.key}"
  }
}
