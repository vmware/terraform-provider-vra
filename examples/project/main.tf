provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_cloud_account_aws" my-cloud-account {
       name = var.cloud_account
}

data "vra_region" "us-east-1-region" {
  cloud_account_id = vra_cloud_account_aws.my-cloud-account.id
  region           = var.region
}

data "vra_zone" "my-zone" {
  name = var.zone
}

resource "vra_project" "my-project" {
  name        = "tf-vra-project"
  description = "my terraform test project"
  zone_assignments {
    zone_id       = vra_zone.my-zone.id
    priority      = 1
    max_instances = 2
  }
}