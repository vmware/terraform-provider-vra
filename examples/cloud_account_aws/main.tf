provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

resource "vra_cloud_account_aws" "this" {
  name        = "tf-vra-cloud-account-aws"
  description = "terraform test cloud account aws"
  access_key  = var.access_key
  secret_key  = var.secret_key
  regions     = ["us-east-1", "us-west-1"]

  tags {
    key   = "foo"
    value = "bar"
  }
}

data "vra_region" "region_1" {
  cloud_account_id = vra_cloud_account_aws.this.id
  region           = "us-east-1"
}

data "vra_region" "region_2" {
  cloud_account_id = vra_cloud_account_aws.this.id
  region           = "us-west-1"
}
