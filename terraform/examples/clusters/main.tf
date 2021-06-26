terraform {
  required_providers {
    tca = {
      version = "0.2"
      source  = "github.com/spyroot/tca"
    }
  }
}

provider "tca" {
  # read from env or provider config
  #tca_username = ""
  #tca_password = ""
  #tca_url = ""
}

data "tca_clusters" "all" {}

output "first_order" {
  value = data.tca_clusters.all
}
