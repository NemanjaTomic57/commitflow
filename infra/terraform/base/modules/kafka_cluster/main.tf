##################################################
# Security Groups
##################################################

resource "aws_security_group" "bastion" {
  name        = "${var.name}-bastion-sg"
  description = "Security group for Bastion"
  vpc_id      = var.vpc_id

  ingress {
    description = "Allow SSH"
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Allow all outbound traffic"
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.name}-bastion-sg"
  }
}

resource "aws_security_group" "kafka" {
  name        = "${var.name}-kafka-sg"
  description = "Security group for Kafka nodes"
  vpc_id      = var.vpc_id

  ingress {
    description     = "Allow SSH from Bastion"
    protocol        = "tcp"
    from_port       = 22
    to_port         = 22
    security_groups = [aws_security_group.bastion.id]
  }

  egress {
    description = "Allow all outbound traffic"
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.name}-kafka-sg"
  }
}

##################################################
# Bastion
##################################################

resource "aws_instance" "bastion" {
  ami                    = var.ami_id
  instance_type          = var.instance_types["bastion"]
  key_name               = var.key_name
  subnet_id              = var.public_subnet_ids[0]
  vpc_security_group_ids = [aws_security_group.bastion.id]

  tags = {
    Name = "${var.name}-ec2-bastion"
  }
}

##################################################
# Kafka nodes
##################################################

locals {
  kafka_nodes = {
    node-1 = var.private_subnet_ids[0]
    node-2 = var.private_subnet_ids[1]
  }
}

resource "aws_instance" "bastion" {
  for_each = local.kafka_nodes

  ami                    = var.ami_id
  instance_type          = var.instance_types["kafka"]
  key_name               = var.key_name
  subnet_id              = each.value
  vpc_security_group_ids = [aws_security_group.kafka.id]

  tags = {
    Name = "${var.name}-ec2-kafka-${each.key}"
  }
}
