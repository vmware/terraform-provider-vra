######################################
# Providers
######################################
provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure = true
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

# Lookup AWS storage profile once it is created, using specific vra_storage_profile_aws data source.
data "vra_storage_profile_aws" "this" {
  id = vra_storage_profile_aws.this.id
}

######################################
# Resources
######################################
# AWS storage profile using generic vra_storage_profile resource. Use 'vra_storage_profile_aws' resource as an alternative.
resource "vra_storage_profile" "this" {
  name                = "aws-with-instance-store"
  description         = "AWS Storage Profile with instance store device type."
  region_id           = data.vra_region.this.id
  default_item        = false
  supports_encryption = false

  disk_properties = {
    deviceType = "ebs"

    // Volume Types: gp2 - General Purpose SSD, io1 - Provisioned IOPS SSD, sc1 - Cold HDD, ST1 - Throughput Optimized HDD, standard - Magnetic
    volumeType = "io1"  // Supported values: gp2, io1, sc1, st1, standard.
    iops       = "1000" // Required only when volumeType is io1.
  }

  tags {
    key   = "foo"
    value = "bar"
  }
}

# AWS storage profile using vra_storage_profile_aws resource.
resource "vra_storage_profile_aws" "this" {
  name                = "aws-with-instance-store-1"
  description         = "AWS Storage Profile with instance store device type."
  region_id           = data.vra_region.this.id
  default_item        = false
  supports_encryption = false

  device_type         = "ebs"

  // Volume Types: gp2 - General Purpose SSD, io1 - Provisioned IOPS SSD, sc1 - Cold HDD, ST1 - Throughput Optimized HDD, standard - Magnetic
  volume_type         = "io1"  // Supported values: gp2, io1, sc1, st1, standard.
  iops                = "1000" // Required only when volumeType is io1.

  tags {
    key   = "foo"
    value = "bar"
  }
}
