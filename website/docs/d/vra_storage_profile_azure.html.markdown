---
layout: "vra"
page_title: "VMware vRealize Automation: vra_storage_profile_azure"
description: |-
  Provides a data lookup for vra_storage_profile_azure.
---

# Resource: vra_storage_profile_azure
## Example Usages
This is an example of how to create a storage profile azure resource.

**Storage profile azure data source by its id:**

```hcl
data "vra_storage_profile_azure" "this" {
  id = vra_storage_profile_azure.this.id
}
```

**Vra storage profile data source filter by external region id:**

```hcl
data "vra_storage_profile_azure" "this" {
  filter = "externalRegionId eq 'foobar'"
}
```

A storage profile azure data source supports the following arguments:

## Argument Reference

* `filter` - (Optional) Filter query string that is supported by vRA multi-cloud IaaS API. Example: `regionId eq '<regionId>' and cloudAccountId eq '<cloudAccountId>'`.

* `id` - (Optional) The id of the image profile instance.

## Attributes Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `data_disk_caching` - Indicates the caching mechanism for additional disk.

* `default_item` - Indicates if this storage profile is a default profile.

* `description` - A human-friendly description.

* `disk_type` -  Indicates the performance tier for the storage type. Premium disks are SSD backed and Standard disks are HDD backed.

* `external_region_id` - The id of the region as seen in the cloud provider for which this profile is defined.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `os_disk_caching` - Indicates the caching mechanism for OS disk. Default policy for OS disks is Read/Write.

* `storage_account_id` - Id of a storage account where in the disk is placed.

* `supports_encryption` - Indicates whether this storage policy should support encryption or not.

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `region_id` - A link to the region that is associated with the storage profile.

* `tags` - A set of tag keys and optional values that were set on this Network Profile.
                      example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
