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

resource "vra_zone" "this" {
  name        = "tf-vra-zone"
  description = "my terraform test cloud zone"
  region_id   = data.vra_region.this.id

  tags {
    key   = "my-tf-key"
    value = "my-tf-value"
  }

  tags {
    key   = "tf-foo"
    value = "tf-bar"
  }
}
