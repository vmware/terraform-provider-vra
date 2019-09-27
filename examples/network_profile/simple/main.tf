provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_cloud_account_aws" "this" {
  name = var.cloud_account
}

data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_aws.this.id
  region           = var.region
}

data "vra_fabric_network" "subnet_1" {
  filter = "name eq '${var.subnet_name_1}' and cloudAccountId eq '${data.vra_cloud_account_aws.this.id}' and externalRegionId eq '${var.region}'"
}

data "vra_fabric_network" "subnet_2" {
  filter = "name eq '${var.subnet_name_2}'"
}

resource "vra_network_profile" "simple" {
  name        = "no-isolation"
  description = "Simple Network Profile with no isolation."
  region_id   = data.vra_region.this.id

  fabric_network_ids = [
    data.vra_fabric_network.subnet_1.id,
    data.vra_fabric_network.subnet_2.id
  ]

  isolation_type = "NONE"

  tags {
    key   = "foo"
    value = "bar"
  }
}
