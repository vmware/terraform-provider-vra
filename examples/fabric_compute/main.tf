provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure
}

resource "vra_fabric_compute" "this" {
  tags {
    key   = "foo"
    value = "bar"
  }
}

