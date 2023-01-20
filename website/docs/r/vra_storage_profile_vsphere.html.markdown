---
layout: "vra"
page_title: "VMware vRealize Automation: vra_storage_profile_vsphere"
description: |-
  Provides a data lookup for vra_storage_profile_vsphere.
---

# Resource: vra_storage_profile_vsphere
## Example Usages
This is an example of how to create a storage profile vsphere resource.

**Vra storage profile vsphere:**

```hcl
# vSphere storage profile using generic vra_storage_profile resource.
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
```

A storage profile vsphere resource supports the following arguments:

## Argument Reference

* `datastore_id` - (Optional) Id of the vSphere Datastore for placing disk and VM.

* `default_item` - (Required) Indicates if this storage profile is a default profile.

* `description` - (Optional) A human-friendly description.

* `disk_mode` - (Optional) Type of mode for the disk.

* `disk_type` - (Optional) Indicates the performance tier for the storage type. Premium disks are SSD backed and Standard disks are HDD backed.

* `limit_iops` - (Optional) The upper bound for the I/O operations per second allocated for each virtual disk.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `provisioning_type` - (Optional) Type of provisioning policy for the disk.

* `region_id` - (Required) The Id of the region that is associated with the storage profile.

* `shares` - (Optional) A specific number of shares assigned to each virtual machine.

* `shares_level` - (Optional) Indicates whether this storage profile supports encryption or not.

* `storage_policy_id` - (Optional) Id of the vSphere Storage Policy to be applied.

* `supports_encryption` - (Optional) Indicates whether this storage policy should support encryption or not.

## Attributes Reference

* `cloud_account_id` - Id of the cloud account this storage profile belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_region_id` - The id of the region as seen in the cloud provider for which this profile is defined.

* `links` - HATEOAS of the entity

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.


* `tags` - A set of tag keys and optional values that were set on this Network Profile.
                      example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
