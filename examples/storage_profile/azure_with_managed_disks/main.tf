######################################
# Providers
######################################
provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = true
}

######################################
# Data sources
######################################
# Lookup Azure cloud account using cloud account name
data "vra_cloud_account_azure" "this" {
  name = var.cloud_account
}

# Lookup Azure region using region name (eastus, etc.)
data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_azure.this.id
  region           = var.region
}

# Lookup Azure storage profile using generic vra_storage_profile data source.
data "vra_storage_profile" "this" {
  id = vra_storage_profile.this.id
}

# Lookup Azure storage profile using specific vra_storage_profile_azure data source.
data "vra_storage_profile_azure" "this" {
  id = vra_storage_profile_azure.this.id
}


######################################
# Resources
######################################
# Azure storage profile using generic vra_storage_profile resource. Use 'vra_storage_profile_azure' resource as an alternative.
resource "vra_storage_profile" "this" {
  name                = "azure-with-managed-disks"
  description         = "Azure Storage Profile with managed disks."
  region_id           = data.vra_region.this.id
  default_item        = false
  supports_encryption = false

  disk_properties = {
    azureDataDiskCaching = "None"         // Supported Values: None, ReadOnly, ReadWrite
    azureManagedDiskType = "Standard_LRS" // Supported Values: Standard_LRS, StandardSSD_LRS, Premium_LRS
    azureOsDiskCaching   = "None"         // Supported Values: None, ReadOnly, ReadWrite
  }

  tags {
    key   = "foo"
    value = "bar"
  }
}

# Azure storage profile using vra_storage_profile_azure resource.
resource "vra_storage_profile_azure" "this" {
  name                = "azure-with-managed-disks-1"
  description         = "Azure Storage Profile with managed disks."
  region_id           = data.vra_region.this.id
  default_item        = false
  supports_encryption = false

  data_disk_caching   = "None"         // Supported Values: None, ReadOnly, ReadWrite
  disk_type           = "Standard_LRS" // Supported Values: Standard_LRS, StandardSSD_LRS, Premium_LRS
  os_disk_caching     = "None"         // Supported Values: None, ReadOnly, ReadWrite

  tags {
    key   = "foo"
    value = "bar"
  }
}
