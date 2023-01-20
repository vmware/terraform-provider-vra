---
layout: "vra"
page_title: "VMware vRealize Automation: vra_storage_profile_vsphere"
description: |-
  Provides a data lookup for vra_storage_profile_vsphere.
---

# Data Source: vra_storage_profile_vsphere
## Example Usages
This is an example of how to create a storage profile vsphere data source.

**Storage profile vsphere data source by its id:**

```hcl
data "vra_storage_profile_vsphere" "this" {
  id = vra_storage_profile_vsphere.this.id
}
```

**Vra storage profile data source filter by external region id:**

```hcl
data "vra_storage_profile_vsphere" "this" {
  filter = "externalRegionId eq 'foobar'"
}
```

A storage profile vsphere data source supports the following arguments:

## Argument Reference

* `filter` - (Optional) Filter query string that is supported by vRA multi-cloud IaaS API. Example: `regionId eq '<regionId>' and cloudAccountId eq '<cloudAccountId>'`.

* `id` - (Optional) The id of the image profile instance.

* `shares_level` - (Optional) Indicates whether this storage profile supports encryption or not.

## Attributes Reference
* `cloud_account_id` - Id of the cloud account this storage profile belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `default_item` - Indicates if this storage profile is a default profile.

* `description` - A human-friendly description.

* `disk_mode` -  Type of mode for the disk.

* `disk_type` -  Indicates the performance tier for the storage type. Premium disks are SSD backed and Standard disks are HDD backed.

* `external_region_id` - The id of the region as seen in the cloud provider for which this profile is defined.

* `limit_iops` - The upper bound for the I/O operations per second allocated for each virtual disk.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `provisioning_type` - Type of provisioning policy for the disk.

* `shares` - A specific number of shares assigned to each virtual machine.

* `supports_encryption` - Indicates whether this storage policy should support encryption or not.

* `tags` - A set of tag keys and optional values that were set on this Network Profile.
                      example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
