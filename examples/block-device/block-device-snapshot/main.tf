provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure = true
}

data "vra_cloud_account_vsphere" "this" {
  name = var.cloud_account
}

data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_vsphere.this.id
  region           = var.region
}

data "vra_project" "this" {
  name = var.project
}

# Lookup vSphere fabric datastore using its name
data "vra_fabric_datastore_vsphere" "this" {
  filter = "name eq '${var.datastore_name}' and externalRegionId eq '${var.region}' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
}

# Lookup vSphere fabric storage policy using its name
data "vra_fabric_storage_policy_vsphere" "this" {
  filter = "name eq '${var.storage_policy_name}' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
}

//To create block-device snapshots, a Storage Profile for vSphere with first class disk type is needed
resource "vra_storage_profile" "this" {
  name         = "vSphere-first-class-disk"
  description  = "vSphere Storage Profile with first class disk."
  region_id    = data.vra_region.this.id
  default_item = true

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

// This block device depends on the storage profile created, and the constraints matches the tags of the storage profile
resource "vra_block_device" "disk1" {
  capacity_in_gb = 10
  name = "terraform_vra_block_device_1"
  project_id = data.vra_project.this.id
  depends_on = [vra_storage_profile.this]
  persistent = true

  constraints {
    mandatory  = true
    expression = "foo:bar"
  }
}

// This block device will be created based on the dataStore and storagePolicy defined in the custom_properties
// If the dataStore and storagePolicy are different from the ones used in the vra_storage_profile,
// then this disk will land into a different datastore
resource "vra_block_device" "disk2" {
  capacity_in_gb = 10
  name = "terraform_vra_block_device_2"
  project_id = data.vra_project.this.id
  custom_properties = {
    "dataStore" = var.datastore_name_b
    "storagePolicy" = var.storage_policy_name_b
  }
  persistent = true
  purge = true
}

resource "vra_block_device_snapshot" "snapshot1" {
  block_device_id = vra_block_device.disk1.id
  description = "terraform fcd snapshot"
}