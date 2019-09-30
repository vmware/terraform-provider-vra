provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_data_collector" "dc" {
  name = var.datacollector
}

resource "vra_cloud_account_nsxt" "this" {
  name        = "tf-nsx-t-account"
  description = "foobar"
  username    = var.username
  password    = var.password
  hostname    = var.hostname
  dc_id        = data.vra_data_collector.dc.id

  accept_self_signed_cert = true

  tags {
    key   = "foo"
    value = "bar"
  }
}
