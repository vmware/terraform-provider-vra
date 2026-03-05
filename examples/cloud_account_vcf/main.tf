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

resource "vra_cloud_account_vcf" "this" {
  name                 = "tf-vra-cloud-account-vcf"
  description          = "tf test vcf cloud account"
  workload_domain_id   = var.workload_domain_id
  workload_domain_name = var.workload_domain_name

  vcenter_hostname = var.vcenter_hostname
  vcenter_password = var.vcenter_password
  vcenter_username = var.vcenter_username

  nsx_hostname = var.nsx_hostname
  nsx_password = var.nsx_password
  nsx_username = var.nsx_username

  dc_id                   = var.data_collector_name != "" ? data.vra_data_collector.this[0].id : "" // Required for vRA Cloud, Optional for vRA on-prem
  sddc_manager_id         = var.sddc_manager_id
  regions                 = var.regions
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
