terraform {
  required_providers {
    tca = {
      version = "0.2"
      source  = "github.com/vmware/tca"
    }
  }
}

variable "vnf_instance_name" {
  type    = string
  default = "test cnf"
}

data "tca_cnfs" "all" {}

# Returns all coffees
output "all_cnfs" {
  value = data.tca_cnfs.all.cnfs
}

# Only returns packer spiced latte
output "cnf" {
  value = {
    for cnf in data.tca_cnfs.all.cnfs:
    cnf.id => cnf
    if cnf.vnf_instance_name == var.vnf_instance_name
  }
}
