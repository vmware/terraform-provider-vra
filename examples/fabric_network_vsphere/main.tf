
provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = "true"
}



resource "vra_fabric_network_vsphere" "simple" {
  cidr            = var.cidr
  default_gateway = var.gateway
  domain          = var.domain
  tags {
    key   = "foo"
    value = "bar"
  }
}

