variable subnet_id {
  description = "Subnet ID"
  type        = list
}

variable subnet_app_id {
  description = "Subnet ID for Application"
  type        = list
}

variable vpc_id {
  description = "VPC ID"
  type        = string
}

variable vpc_app_id {
  description = "VPC ID for Application"
  type        = string
}

variable projectname {
  type = string
}

variable image {
  type = string
  default = "docker.io/ubuntu"
}

variable execrole {
  type = string
  default = "arn:aws:iam::225563488370:role/ecsTaskExecutionRole"
}

variable rep_url {
  type = string
}