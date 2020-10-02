######################################
# Providers
######################################
provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

######################################
# Data sources
######################################
# Lookup AWS cloud account using cloud account name
data "vra_cloud_account_aws" "this" {
  name = var.cloud_account
}

# Lookup AWS region using region name (us-east-1, etc.)
data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_aws.this.id
  region           = var.region
}

# Lookup AWS storage profile once it is created, using generic vra_storage_profile data source.
data "vra_storage_profile" "this" {
  id = vra_storage_profile.this.id
}

######################################
# Resources
######################################
# AWS storage profile using generic vra_storage_profile resource. Use 'vra_storage_profile_aws' resource as an alternative.
resource "vra_storage_profile" "this" {
  name         = "aws-with-instance-store"
  description  = "AWS Storage Profile with instance store device type."
  region_id    = data.vra_region.this.id
  default_item = false

  disk_properties = {
    deviceType = "instance-store"
  }

  tags {
    key   = "foo"
    value = "bar"
  }
}
