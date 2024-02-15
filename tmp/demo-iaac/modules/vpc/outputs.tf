output "vpc_main" {
  value = aws_vpc.project-vpc-pub.id
}

output "vpc_app" {
  value = aws_vpc.project-vpc-app.id
}

output "subnet_main" {
  value = aws_subnet.project-subnet-pub[*]
}

output "subnet_db" {
  value = aws_subnet.project-subnet-db[*]
}

output "subnet_app" {
  value = aws_subnet.project-subnet-app[*]
}

output "sg_main" {
  value = aws_security_group.project.id
}
