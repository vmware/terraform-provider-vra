provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_cloud_account_aws" aws-cloud-account {
       name = var.cloud_account
}

data "vra_region" "us-west-1-region" {
  cloud_account_id = data.vra_cloud_account_aws.aws-cloud-account.id
  region           = var.region
}

data "vra_fabric_network" "my-subnet-1" {
  filter = "name eq '${var.subnet-name-1}' and cloudAccountId eq '${data.vra_cloud_account_aws.aws-cloud-account.id}' and externalRegionId eq '${var.region}'"
}

data "vra_fabric_network" "my-subnet-2" {
  filter = "name eq '${var.subnet-name-2}'"
}

resource "vra_network_profile" "my-network-profile" {
  name = "my-vra-network-profile"
  description = "my network profile"
  region_id = data.vra_region.us-west-1-region.id
  fabric_network_ids = [ data.vra_fabric_network.my-subnet-1.id, data.vra_fabric_network.my-subnet-2.id ]
  isolation_type = "NONE"
}
