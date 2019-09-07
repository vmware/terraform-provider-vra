provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

resource "vra_cloud_account_gcp" "this" {
  name           = "tf-vra-cloud-account-gcp"
  description    = "terraform test cloud account gcp"
  client_email   = var.client_email
  private_key_id = var.private_key_id
  private_key    = var.private_key
  project_id     = var.project_id
  regions        = ["us-west1", "us-west2"]

  tags {
    key   = "foo"
    value = "bar"
  }
}

data "vra_region" "region_1" {
  cloud_account_id = vra_cloud_account_gcp.this.id
  region           = "us-west1"
}

data "vra_region" "region_2" {
  cloud_account_id = vra_cloud_account_gcp.this.id
  region           = "us-west2"
}
