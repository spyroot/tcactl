terraform {
  required_providers {
    tca = {
      version = "0.2"
      source = "github.com/spyroot/tca"
    }
  }
}

provider "tca" {
  # read from env or provider config
  #tca_username = ""
  #tca_password = ""
  #tca_url = ""
}

data "tca_templates" "all" {}

# All templates
//output "all" {
//  value = data.tca_templates.all
//}

# template type workload
variable "cluster_type_workload" {
  type = string
  default = "WORKLOAD"
}

# Only returns workload cluster templates
//output "workload_templates" {
//  value = {
//  for c in data.tca_templates.all.templates:
//  c.id => c
//  if c.cluster_type == var.cluster_type_workload
//  }
//}

output "workload_templates" {
  value = {
  for c in data.tca_templates.all.templates:
  c.id => c
  if c.cluster_type == var.cluster_type_workload
  }
}