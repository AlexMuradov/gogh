resource "aws_db_subnet_group" "project" {
  name       = "main"
  subnet_ids = [for subnet in var.subnet_id : subnet.id]
}

resource "aws_db_instance" "project" {
  allocated_storage    = 10
  db_name              = "mydb"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t3.micro"
  username             = var.db_user
  password             = var.db_password
  parameter_group_name = "default.mysql5.7"
  skip_final_snapshot  = true
  db_subnet_group_name = aws_db_subnet_group.project.name
}