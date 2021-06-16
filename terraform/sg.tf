resource "aws_security_group" "compliant" {
  name_prefix = "compliant_sg_"
  description = "compliant securitygroup"
  vpc_id      = var.vpc_id
  dynamic "ingress" {
    for_each = var.compliant_ingress
    content {
      description = ingress.key
      from_port   = ingress.value.port_range[0]
      to_port     = ingress.value.port_range[1]
      protocol    = ingress.value.protocol
      cidr_blocks = ingress.value.source_cidrs
    }
  }
  tags = {
    "Name" = "good_sg"
    "Interviewee" : "RyanODonnell"
  }
}


resource "aws_security_group" "non_compliant" {
  name_prefix = "non_compliant_sg_"
  description = "non_compliant securitygroup"
  vpc_id      = var.vpc_id
  dynamic "ingress" {
    for_each = var.non_compliant_ingress
    content {
      description = ingress.key
      from_port   = ingress.value.port_range[0]
      to_port     = ingress.value.port_range[1]
      protocol    = ingress.value.protocol
      cidr_blocks = ingress.value.source_cidrs
    }
  }
  tags = {
    "Name" = "bad_sg"
    "Interviewee" : "RyanODonnell"
  }
}


variable "vpc_id" {

}



variable "non_compliant_ingress" {
  default = {
    ssh = {
      protocol     = "tcp"
      source_cidrs = ["0.0.0.0/0", "10.0.0.0/24"]
      port_range   = [22, 22]
    }
    sql = {
      protocol     = "tcp"
      source_cidrs = ["0.0.0.0/0", "10.0.0.0/24"]
      port_range   = [3306, 3306]
    }
    rdp = {
      protocol     = "tcp"
      source_cidrs = ["0.0.0.0/0", "10.0.0.0/24"]
      port_range   = [3389, 3389]
    }
  }
}



variable "compliant_ingress" {
  default = {
    https = {
      protocol     = "tcp"
      source_cidrs = ["0.0.0.0/0", "10.0.0.0/24"]
      port_range   = [443, 443]
    }
    http = {
      protocol     = "tcp"
      source_cidrs = ["0.0.0.0/0", "10.0.0.0/24"]
      port_range   = [80, 80]
    }
    http_8080 = {
      protocol     = "tcp"
      source_cidrs = ["0.0.0.0/0", "10.0.0.0/24"]
      port_range   = [8080, 8080]
    }
  }
}
