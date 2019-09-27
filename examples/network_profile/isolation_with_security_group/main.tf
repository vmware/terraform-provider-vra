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

resource "vra_network_profile" "security_group_isolation" {
  name        = "isolation-with-security-group"
  description = "On-demand security groups are created for outbound and private networks."
  region_id   = data.vra_region.this.id

  fabric_network_ids = [
    data.vra_fabric_network.subnet.id
  ]

  isolation_type = "SECURITY_GROUP"

  tags {
    key   = "foo"
    value = "bar"
  }
}
