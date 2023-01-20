---
layout: "vra"
page_title: "VMware vRealize Automation: vra_storage_profile_azure"
description: |-
  Provides a data lookup for vra_storage_profile_azure.
---

# Resource: vra_storage_profile_azure
## Example Usages
This is an example of how to create a storage profile azure resource.

**Vra storage profile azure:**

```hcl
# Azure storage profile using vra_storage_profile_azure resource with managed disk.
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

# Azure storage profile using vra_storage_profile_azure resource with unmanaged disk.
resource "vra_storage_profile_azure" "this" {
  name                = "azure-with-unmanaged-disks"
  description         = "Azure Storage Profile with unmanaged disks."
  region_id           = data.vra_region.this.id
  default_item        = false
  supports_encryption = false

  data_disk_caching   = "None" // Supported Values: None, ReadOnly, ReadWrite
  os_disk_caching     = "None" // Supported Values: None, ReadOnly, ReadWrite

  tags {
    key   = "foo"
    value = "bar"
  }
}
```

A storage profile azure resource supports the following arguments:

## Argument Reference

* `data_disk_caching` - (Optional) Indicates the caching mechanism for additional disk.

* `default_item` - (Required) Indicates if this storage profile is a default profile.

* `description` - (Optional) A human-friendly description.

* `disk_type` - (Optional) Indicates the performance tier for the storage type. Premium disks are SSD backed and Standard disks are HDD backed.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `os_disk_caching` - (Optional) Indicates the caching mechanism for OS disk. Default policy for OS disks is Read/Write.

* `region_id` - (Required) A link to the region that is associated with the storage profile.

* `storage_account_id` - (Optional) Id of a storage account where in the disk is placed.

* `supports_encryption` - (Optional) Indicates whether this storage policy should support encryption or not.

## Attributes Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_region_id` - The id of the region as seen in the cloud provider for which this profile is defined.

* `links` - HATEOAS of the entity

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` - A set of tag keys and optional values that were set on this Network Profile.
                      example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
