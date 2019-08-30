provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_data_collector" "dc" {
    name = var.datacollector
}

data "vra_region_enumeration" "dc_regions" {
  username    = var.username
  password    = var.password
  hostname    = var.hostname
  dcid        = data.vra_data_collector.dc.id
}

resource "vra_cloud_account_vsphere" "my_vsphere_account" {
  name        = "my-vsphere-account"
  description = "foobar"
  username    = var.username
  password    = var.password
  hostname    = var.hostname
  dcid        = data.vra_data_collector.dc.id

  regions                 = data.vra_region_enumeration.dc_regions.regions
  accept_self_signed_cert = true
  tags {
    key   = "foo"
    value = "bar"
  }
}
