resource "aws_vpc" "project-vpc-pub" {
  cidr_block = "10.${var.vpc_pub[0]["iprange"]}.0.0/16"
  tags = {
    Name = "${var.projectname} ${var.vpc_pub[0]["desc"]}"
  }
}

resource "aws_vpc" "project-vpc-db" {
  cidr_block = "10.${var.vpc_db[0]["iprange"]}.0.0/16"

  tags = {
    Name = "${var.projectname} ${var.vpc_db[0]["desc"]}"
  }
}

resource "aws_vpc" "project-vpc-app" {
  cidr_block = "10.${var.vpc_app[0]["iprange"]}.0.0/16"

  tags = {
    Name = "${var.projectname} ${var.vpc_app[0]["desc"]}"
  }
}

resource "aws_internet_gateway" "project-igw" {
  vpc_id = aws_vpc.project-vpc-pub.id

  tags = {
    Name = var.igw_pub[0]["desc"]
  }
}

resource "aws_route_table" "project-rt" {
  vpc_id = aws_vpc.project-vpc-pub.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.project-igw.id
  }

  tags = {
    Name = var.igw_pub[0]["desc"]
  }
}

resource "aws_subnet" "project-subnet-pub" {
  count = length(var.subnets)
  vpc_id     = aws_vpc.project-vpc-pub.id
  cidr_block = "10.${var.vpc_pub[0]["iprange"]}.${count.index + 1}.0/24"
  availability_zone = var.subnets[count.index]["az"]

  tags = {
    Name = "${var.projectname} subnet ${var.subnets[count.index]["az"]}"
  }
}

resource "aws_subnet" "project-subnet-db" {
  count = length(var.subnets)
  vpc_id     = aws_vpc.project-vpc-db.id
  cidr_block = "10.${var.vpc_db[0]["iprange"]}.${count.index + 1}.0/24"
  availability_zone = var.subnets[count.index]["az"]

  tags = {
    Name = "${var.projectname} subnet ${var.subnets[count.index]["az"]}"
  }
}

resource "aws_subnet" "project-subnet-app" {
  count = length(var.subnets)
  vpc_id     = aws_vpc.project-vpc-app.id
  cidr_block = "10.${var.vpc_app[0]["iprange"]}.${count.index + 1}.0/24"
  availability_zone = var.subnets[count.index]["az"]

  tags = {
    Name = "${var.projectname} subnet ${var.subnets[count.index]["az"]}"
  }
}

resource "aws_security_group" "project" {
  vpc_id      = aws_vpc.project-vpc-pub.id
  ingress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}