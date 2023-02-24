provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure // false for vRA Cloud and true for vRA 8.0
}

# Required for vRA Cloud, Optional for vRA 8.X
data "vra_data_collector" "dc" {
  count = var.datacollector != "" ? 1 : 0
  name  = var.datacollector
}

resource "vra_cloud_account_nsxt" "this" {
  name        = "tf-nsx-t-account"
  description = "foobar"
  username    = var.nsxt_username
  password    = var.nsxt_password
  hostname    = var.nsxt_hostname
  dc_id       = var.datacollector != "" ? data.vra_data_collector.dc[0].id : "" // Required for vRA Cloud, Optional for vRA 8.X

  accept_self_signed_cert = true

  tags {
    key   = "foo"
    value = "bar"
  }
}

data "vra_region_enumeration_vsphere" "this" {
  username                = var.username
  password                = var.password
  hostname                = var.hostname
  dc_id                   = var.datacollector != "" ? data.vra_data_collector.dc[0].id : "" // Required for vRA Cloud, Optional for vRA 8.X
  accept_self_signed_cert = true
}

resource "vra_cloud_account_vsphere" "this" {
  name        = "tf-vsphere-account"
  description = "foobar"
  username    = var.username
  password    = var.password
  hostname    = var.hostname
  dc_id       = var.datacollector != "" ? data.vra_data_collector.dc[0].id : "" // Required for vRA Cloud, Optional for vRA 8.X

  regions                      = data.vra_region_enumeration_vsphere.this.regions
  associated_cloud_account_ids = [vra_cloud_account_nsxt.this.id]
  accept_self_signed_cert      = true

  tags {
    key   = "foo"
    value = "bar"
  }
}
