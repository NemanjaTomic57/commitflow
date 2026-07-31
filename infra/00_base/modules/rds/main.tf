resource "aws_db_subnet_group" "this" {
  name       = var.name
  subnet_ids = values(var.private_subnet_ids)

  tags = {
    Name = "${var.name}-db-subnet-group"
  }
}

resource "aws_db_instance" "this" {
  db_name                     = var.name
  engine                      = "postgres"
  engine_version              = "18.4"
  username                    = "postgres"
  manage_master_user_password = true
  instance_class              = var.instance_class
  storage_type                = var.storage_type
  allocated_storage           = var.allocated_storage
  db_subnet_group_name        = aws_db_subnet_group.this.name
  vpc_security_group_ids      = [var.db_security_group_id]
  skip_final_snapshot         = true

  tags = {
    Name = "${var.name}-db"
  }
}
