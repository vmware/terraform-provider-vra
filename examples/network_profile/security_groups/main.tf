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

data "vra_fabric_network" "subnet" {
  filter = "name eq '${var.subnet_name}' and cloudAccountId eq '${data.vra_cloud_account_aws.this.id}' and externalRegionId eq '${var.region}'"
}

data "vra_security_group" "this" {
  filter = "name eq '${var.security_group_name}' and cloudAccountId eq '${data.vra_cloud_account_aws.this.id}' and externalRegionId eq '${var.region}'"
}

resource "vra_network_profile" "firewall_rules" {
  name        = "network-profile-with-firewall-rules"
  description = "Firewall rules are added to all machines provisioned."
  region_id   = data.vra_region.this.id

  fabric_network_ids = [
    data.vra_fabric_network.subnet.id
  ]

  security_group_ids = [
    data.vra_security_group.this.id
  ]

  tags {
    key   = "foo"
    value = "bar"
  }
}
