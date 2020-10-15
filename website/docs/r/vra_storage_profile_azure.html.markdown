---
layout: "vra"
page_title: "VMware vRealize Automation: vra_storage_profile_azure"
sidebar_current: "docs-vra-resource-vra-storage-profile_azure"
description: |-
  Provides a data lookup for vra_storage_profile_azure.
---

# Resource: vra_storage_profile_azure
## Example Usages
This is an example of how to create a storage profile azure resource.

**Vra storage profile azure:**

```hcl
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
```

A storage profile azure resource supports the following arguments:

## Required arguments

* `name` - A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `default_item` - Indicates if this storage profile is a default profile.

* `region_id` - A link to the region that is associated with the storage profile.

## Optional arguments

* `data_disk_caching` - Indicates the caching mechanism for additional disk.

* `description` - A human-friendly description.

* `disk_type` -  Indicates the performance tier for the storage type. Premium disks are SSD backed and Standard disks are HDD backed.

* `os_disk_caching` - Indicates the caching mechanism for OS disk. Default policy for OS disks is Read/Write.

* `storage_account_id` - Id of a storage account where in the disk is placed.

* `supports_encryption` - Indicates whether this storage policy should support encryption or not.

## Imported attributes
* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_region_id` - The id of the region as seen in the cloud provider for which this profile is defined.

* `links` - HATEOAS of the entity

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` - A set of tag keys and optional values that were set on this Network Profile.
                      example:[ { "key" : "ownedBy", "value": "Rainpole" } ]

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
