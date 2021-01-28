terraform {
  required_providers {
    null = {
      source = "hashicorp/null"
    }
    vra = {
      source = "vmware/vra"
    }
  }
  required_version = ">= 0.13"
}
