provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure // false for vRA Cloud and true for vRA 8.0
}

// Required for vRA Cloud, Optional for vRA 8.0
data "vra_data_collector" "dc" {
  count = var.datacollector != "" ? 1 : 0
  name  = var.datacollector
}

data "vra_region_enumeration" "dc_regions" {
  username = var.username
  password = var.password
  hostname = var.hostname
  dcid     = var.datacollector != "" ? data.vra_data_collector.dc[0].id : "" // Required for vRA Cloud, Optional for vRA 8.0
}

resource "vra_cloud_account_vsphere" "this" {
  name        = "tf-vsphere-account"
  description = "foobar"
  username    = var.username
  password    = var.password
  hostname    = var.hostname
  dcid        = var.datacollector != "" ? data.vra_data_collector.dc[0].id : "" // Required for vRA Cloud, Optional for vRA 8.0

  regions                 = data.vra_region_enumeration.dc_regions.regions
  accept_self_signed_cert = true

  tags {
    key   = "foo"
    value = "bar"
  }
}
