terraform {
  required_providers {
    tca = {
      version = "0.2"
      source  = "github.com/vmware/tca"
    }
  }
}

provider "tca" {}

module "psl" {
  source   = "./cnfs"
  cnf_name = "test"
}

output "psl" {
  value = module.psl.cnf
}
