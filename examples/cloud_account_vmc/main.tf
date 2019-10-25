provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure // false for vRA Cloud and true for vRA 8.0
}

// Required for vRA Cloud, Optional for vRA 8.0
data "vra_data_collector" "this" {
  count = var.data_collector_name != "" ? 1 : 0
  name  = var.data_collector_name
}

data "vra_region_enumeration" "this" {
  hostname = var.vcenter_hostname
  password = var.vcenter_password
  username = var.vcenter_username
  dcid     = var.data_collector_name != "" ? data.vra_data_collector.this[0].id : "" // Required for vRA Cloud, Optional for vRA 8.0
}

resource "vra_cloud_account_vmc" "this" {
  name        = "tf-vra-cloud-account-vmc"
  description = "tf test vmc cloud account"

  api_token = var.api_token
  sddc_name = var.sddc_name

  vcenter_hostname = var.vcenter_hostname
  vcenter_password = var.vcenter_password
  vcenter_username = var.vcenter_username
  nsx_hostname     = var.nsx_hostname
  dc_id            = var.data_collector_name != "" ? data.vra_data_collector.this[0].id : "" // Required for vRA Cloud, Optional for vRA 8.0

  regions                 = data.vra_region_enumeration.this.regions
  accept_self_signed_cert = true

  tags {
    key   = "foo"
    value = "bar"
  }

  tags {
    key   = "where"
    value = "waldo"
  }
}

data "vra_region" "region_1" {
  cloud_account_id = vra_cloud_account_vmc.this.id
  region           = data.vra_region_enumeration.this.regions[0]
}
