provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_cloud_account_aws" this {
  name = var.cloud_account
}

data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_aws.this.id
  region           = var.region
}

data "vra_zone" "this" {
  name = var.zone
}

resource "vra_project" "this" {
  name        = "tf-vra-project"
  description = "terraform test project"

  zone_assignments {
    zone_id       = data.vra_zone.this.id
    priority      = 1
    max_instances = 2
  }
}
