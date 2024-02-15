resource "aws_ecr_repository" "project" {
  name                 = "${var.projectname}-registry"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = false
  }
}