provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

resource "vra_cloud_account_gcp" "my-cloud-account" {
  name            = "tf-vra-cloud-account-gcp"
  description     = "my terraform test cloud account gcp"
  client_email    = var.client_email
  private_key_id  = var.private_key_id
  private_key     = var.private_key
  project_id      = var.project_id
  regions         = ["us-west1", "us-west2"]
}

data "vra_region" "us-east-1-region" {
  cloud_account_id = vra_cloud_account_gcp.my-cloud-account.id
  region           = "us-west1"
}

data "vra_region" "us-west-1-region" {
  cloud_account_id = vra_cloud_account_gcp.my-cloud-account.id
  region           = "us-west2"
}
