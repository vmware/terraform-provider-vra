provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

resource "vra_cloud_account_zure" "this" {
  name        = "tf-vra-cloud-account-azure"
  description = "terraform test cloud account azure"

  subscription_id = var.subscription_id
  tenant_id       = var.tenant_id
  application_id  = var.application_id
  application_key = var.application_key

  regions = ["eastus", "centralus"]

  tags {
    key   = "foo"
    value = "bar"
  }
  tags {
    key   = "where"
    value = "waldo"
  }
}

data "vra_region" "eastus" {
  cloud_account_id = vra_cloud_account_azure.this.id
  region           = "eastus"
}

data "vra_region" "centralus" {
  cloud_account_id = vra_cloud_account_azure.this.id
  region           = "centralus"
}
