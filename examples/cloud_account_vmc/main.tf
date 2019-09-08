provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_data_collector" "this" {
  name = var.data_collector_name
}

data "vra_region_enumeration" "this" {
  hostname = var.vcenter_hostname
  password = var.vcenter_password
  username = var.vcenter_username
  dcid     = data.vra_data_collector.this.id
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
  dc_id            = data.vra_data_collector.this.id

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
