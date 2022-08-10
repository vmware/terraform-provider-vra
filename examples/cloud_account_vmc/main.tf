provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure // false for vRA Cloud and true for vRA on-prem
}

// Required for vRA Cloud, Optional for vRA on-prem
data "vra_data_collector" "this" {
  count = var.data_collector_name != "" ? 1 : 0
  name  = var.data_collector_name
}

data "vra_region_enumeration_vmc" "this" {
  api_token = var.api_token
  sddc_name = var.sddc_name
  nsx_hostname = var.nsx_hostname

  vcenter_hostname = var.vcenter_hostname
  vcenter_password = var.vcenter_password
  vcenter_username = var.vcenter_username
  dc_id            = var.data_collector_name != "" ? data.vra_data_collector.this[0].id : "" // Required for vRA Cloud, Optional for vRA on-prem
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
  dc_id            = var.data_collector_name != "" ? data.vra_data_collector.this[0].id : "" // Required for vRA Cloud, Optional for vRA on-prem

  regions                 = data.vra_region_enumeration_vmc.this.regions
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
  region           = tolist(data.vra_region_enumeration.this.regions)[0]
}
