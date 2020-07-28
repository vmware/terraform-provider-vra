provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure
}

data "vra_data_collector" "dc" {
  count = var.cloud_proxy != "" ? 1 : 0
  name  = var.cloud_proxy
}

resource "vra_cloud_account_nsxv" "this" {
  name        = "tf-nsx-v-account"
  description = "foobar"
  username    = var.username
  password    = var.password
  hostname    = var.hostname
  dc_id       = var.cloud_proxy != "" ? data.vra_data_collector.dc[0].id : ""

  accept_self_signed_cert = true

  tags {
    key   = "foo"
    value = "bar"
  }
}
