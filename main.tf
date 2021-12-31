terraform {

  backend "http" {
    # address omitted, must be configured via init
    update_method = "PUT"
  }

  required_providers {
    oci = {
      source  = "hashicorp/oci"
      version = "~> 4.57"
    }
  }
}

provider "oci" {
  auth                = "SecurityToken"
  config_file_profile = "DEFAULT"
  region              = var.region
}
