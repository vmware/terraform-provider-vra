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
# Lookup vSphere cloud account using cloud account name
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

# Lookup vSphere storage profile using generic vra_storage_profile data source.
data "vra_storage_profile" "this" {
  id = vra_storage_profile.this.id
}

######################################
# Resources
######################################
# vSphere storage profile using generic vra_storage_profile resource.
resource "vra_storage_profile" "this" {
  name         = "vSphere-first-class-disk"
  description  = "vSphere Storage Profile with first class disk."
  region_id    = data.vra_region.this.id
  default_item = false

  disk_properties = {
    diskType         = "firstClass"
    provisioningType = "thin" // Supported values: "thin", "thick", "eagerZeroedThick"
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

data "vra_storage_profile_vsphere" "this" {
  id = vra_storage_profile_vsphere.this.id
}

# vSphere storage profile using specific vra_storage_profile_vsphere resource.
resource "vra_storage_profile_vsphere" "this" {
  name = "vra_storage_profile_vsphere resource - FCD"
  description = "vSphere Storage Profile with FCD disk."
  region_id = data.vra_region.this.id
  default_item = false
  disk_type = "firstClass"

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
