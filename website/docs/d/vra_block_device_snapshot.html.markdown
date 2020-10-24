---layout: "vra"
page_title: "VMware vRealize Automation: vra_block_device_snapshots"
description: |-
  Provides a data lookup for vra_block_device_snapshots.
---

# Data Source: vra_block_device_snapshots
## Example Usages

This is an example of how to read a block device snapshots data source.

**Block device snapshots data source by its id:**
```hcl

data "vra_block_device_snapshot" "snapshot" {
  block_device_id = var.block_device_id
}

```
## Argument Reference

* `var.block_device_id` - (Required) The id of the existing block device.

## Attributes Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `description` - A human-friendly description.

* `is_current` - Indicates whether this snapshot is the current snapshot on the block-device.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `update_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
