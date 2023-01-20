---
layout: "vra"
page_title: "VMware vRealize Automation: vra_storage_profile"
description: |-
  Provides a data lookup for vra_storage_profile.
---

# Resource: vra_storage_profile
## Example Usages
This is an example of how to create a storage profile resource.

**Vra storage profile:**

```hcl
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

# AWS storage profile using generic vra_storage_profile resource.
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

# Azure storage profile using generic vra_storage_profile resource.
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
```

A storage profile resource supports the following arguments:

## Argument Reference

* `description` - (Optional) A human-friendly description.

* `disk_properties` - (Optional) Map of storage properties that are to be applied on disk while provisioning.

* `disk_target_properties` - (Optional) Map of storage placements to know where the disk is provisioned.

* `region_id` - (Required) The id of the region for which this profile is defined as in vRealize Automation(vRA).

* `supports_encryption` - (Optional) Indicates whether this storage profile supports encryption or not.

## Attributes Reference

* `cloud_account_id` - Id of the cloud account this storage profile belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_region_id` - The id of the region as seen in the cloud provider for which this profile is defined.

* `default_item` - Indicates if this storage profile is a default profile.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` - A set of tag keys and optional values that were set on this Network Profile.
           example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
