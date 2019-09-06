provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_cloud_account_aws" aws_cloud_account {
       name = var.cloud_account
}

data "vra_region" "example_region" {
  cloud_account_id = data.vra_cloud_account_aws.aws_cloud_account.id
  region           = var.region
}

data "vra_fabric_network" "example_subnet_1" {
  filter = "name eq '${var.subnet_name_1}' and cloudAccountId eq '${data.vra_cloud_account_aws.aws_cloud_account.id}' and externalRegionId eq '${var.region}'"
}

data "vra_fabric_network" "example_subnet_2" {
  filter = "name eq '${var.subnet_name_2}'"
}

resource "vra_network_profile" "example_network_profile" {
  name = "example_network_profile"
  description = "network profile description"
  region_id = data.vra_region.example_region.id
  fabric_network_ids = [ data.vra_fabric_network.example_subnet_1.id, data.vra_fabric_network.example_subnet_2.id ]
  isolation_type = "NONE"
}
