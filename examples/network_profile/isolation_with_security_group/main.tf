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

data "vra_fabric_network" "example_subnet" {
  filter = "name eq '${var.subnet_name}' and cloudAccountId eq '${data.vra_cloud_account_aws.aws_cloud_account.id}' and externalRegionId eq '${var.region}'"
}

resource "vra_network_profile" "example_network_profile" {
  name = "vra_network_profile_example"
  description = "network profile description"
  region_id = data.vra_region.example_region.id
  fabric_network_ids = [ data.vra_fabric_network.example_subnet.id ]
  isolation_type = "SECURITY_GROUP"
}
