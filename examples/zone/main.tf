provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_cloud_account_aws" my-cloud-account {
       name = "tf-vra-cloud-account-aws"
}

data "vra_region" "us-east-1-region" {
  cloud_account_id = vra_cloud_account_aws.my-cloud-account.id
  region           = "us-east-1"
}

resource "vra_zone" "my-zone" {
  name        = "tf-vra-zone"
  description = "my terraform test cloud zone"
  region_id   = data.vra_region.us-east-1-region.id
  tags {
    key   = "my-tf-key"
    value = "my-tf-value"
  }
  tags {
    key   = "tf-foo"
    value = "tf-bar"
  }
}