provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

resource "vra_cloud_account_aws" "my-cloud-account" {
  name        = "tf-vra-cloud-account-aws"
  description = "my terraform test cloud account aws"
  access_key  = var.access_key
  secret_key  = var.secret_key
  regions     = ["us-east-1", "us-west-1"]
}

data "vra_region" "us-east-1-region" {
  cloud_account_id = vra_cloud_account_aws.my-cloud-account.id
  region           = "us-east-1"
}

data "vra_region" "us-west-1-region" {
  cloud_account_id = vra_cloud_account_aws.my-cloud-account.id
  region           = "us-west-1"
}
