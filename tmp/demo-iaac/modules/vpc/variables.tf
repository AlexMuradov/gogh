variable "subnets" {
  type = list
  default = [
      { 
        cidr = "1",
        az = "eu-central-1a"
      },
      {
        cidr = "2",
        az = "eu-central-1b"
      },
      {
        cidr = "3",
        az = "eu-central-1c"
      }
    ]
}

variable "vpc_pub" {
  type = list
  default = [
      { 
        desc = "vpc pub",
        iprange = "10"
      }
    ]
}

variable "vpc_db" {
  type = list
  default = [
      { 
        desc = "vpc db",
        iprange = "11"
      }
    ]
}

variable "vpc_app" {
  type = list
  default = [
      { 
        desc = "vpc app",
        iprange = "12"
      }
    ]
}


variable "igw_pub" {
  type = list
  default = [
      { 
        desc = "internet gateway"
      }
    ]
}

variable projectname {
  type = string
}