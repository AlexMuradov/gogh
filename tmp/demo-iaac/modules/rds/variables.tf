variable subnet_id {
  description = "Subnet ID"
  type        = list
}

variable projectname {
  type = string
}

variable db_user {
    type = string
    default = "root"
}

variable db_password {
    type = string
    default = "rootroot"
}