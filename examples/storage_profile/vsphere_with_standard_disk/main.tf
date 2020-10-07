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
data "vra_cloud_account_vsphere" "this" {
  name = var.cloud_account
}

# Lookup vSphere region using region/datacenter id (moref from vCenter, usually in the format of Datacenter:datacenter-2)
data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_vsphere.this.id
  region           = var.region
}

# Lookup vSphere fabric datastore using its name
data "vra_fabric_datastore_vsphere" "this" {
  filter = "name eq '${var.datastore_name}' and externalRegionId eq '${var.region}' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
}

# Lookup vSphere fabric storage policy using its name
data "vra_fabric_storage_policy_vsphere" "this" {
  filter = "name eq '${var.storage_policy_name}' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
}

# Lookup vSphere storage profile once it is created, using generic vra_storage_profile data source.
data "vra_storage_profile" "this" {
  id = vra_storage_profile.this.id
}

######################################
# Resources
######################################
# vSphere storage profile using generic vra_storage_profile resource.
resource "vra_storage_profile" "this" {
  name         = "vSphere-standard-independent-non-persistent-disk"
  description  = "vSphere Storage Profile with standard independent non-persistent disk."
  region_id    = data.vra_region.this.id
  default_item = false

  disk_properties = {
    independent      = "true"
    persistent       = "false"
    limitIops        = "2000"
    provisioningType = "eagerZeroedThick" // Supported values: "thin", "thick", "eagerZeroedThick"
    sharesLevel      = "custom"           // Supported values: "low", "normal", "high", "custom"
    shares           = "1500"             // Required only when sharesLevel is "custom".
  }

  disk_target_properties = {
    datastoreId     = data.vra_fabric_datastore_vsphere.this.id
    storagePolicyId = data.vra_fabric_storage_policy_vsphere.this.id // Remove it if datastore default storage policy needs to be selected.
  }

  tags {
    key   = "foo"
    value = "bar"
  }
}

# vSphere storage profile using specific vra_storage_profile_vsphere resource.
resource "vra_storage_profile_vsphere" "this" {
  name = "vra_storage_profile_vsphere resource - standard"
  description = "vSphere Storage Profile with standard disk."
  region_id = data.vra_region.this.id
  default_item = false
  disk_type = "standard"

  provisioning_type = "thin"
  // Supported values: "thin", "thick", "eagerZeroedThick"

  datastore_id = data.vra_fabric_datastore_vsphere.this.id
  storage_policy_id = data.vra_fabric_storage_policy_vsphere.this.id
  // Remove it if datastore default storage policy needs to be selected.

  tags {
    key = "foo"
    value = "bar"
  }
}
