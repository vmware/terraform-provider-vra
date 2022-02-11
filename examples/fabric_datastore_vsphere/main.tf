provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure
}

resource "vra_fabric_datastore_vsphere" "this" {
  tags {
    key   = "foo"
    value = "bar"
  }
}

