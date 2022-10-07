provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}


data "vra_cloud_account_aws" "this" {
  name = var.cloud_account
}

data "vra_fabric_network" "subnet" {
	filter = "name eq '${var.subnet_name}' and cloudAccountId eq '${data.vra_cloud_account_aws.this.id}'"
}

resource "vra_network_ip_range" "this" {
  name               = "example-ip-range"
  description        = "Internal Network IP Range Example"
  start_ip_address   = var.start_ip
  end_ip_address     = var.end_ip
  ip_version         = var.ip_version
  fabric_network_ids = [data.vra_fabric_network.subnet.id]

  tags {
    key   = "foo"
    value = "bar"
  }
}
